package functionrecorder

import (
	"fmt"
	"reflect"
	"sync"
	"unicode"
	"unsafe"
)

type (
	valueTree struct {
		root *node
	}

	node struct {
		name, dataType                      string
		dataKind                            reflect.Kind
		data                                interface{}
		associatedNode                      *node // This is for maps
		isSliceMember, isMapKey, isMapValue bool
		children                            []*node
	}
)

func (n *node) append(child *node) *node {
	if len(n.children) == 0 {
		n.children = []*node{child}
	} else {
		n.children = append(n.children, child)
	}

	return child
}

func buildTree(root interface{}, ch chan<- *valueTree) {
	var tree valueTree

	if root != nil {
		rootValue := reflect.ValueOf(root)
		tree = valueTree{
			root: &node{dataType: rootValue.Type().Name(), dataKind: rootValue.Kind(), data: root},
		}

		handleValue(tree.root)
	}

	ch <- &tree
}

func handleValue(n *node) {
	switch reflect.ValueOf(n.data).Kind() {
	case reflect.Array, reflect.Slice:
		handleSlice(n)
	case reflect.Struct:
		handleStruct(n)
	case reflect.Ptr:
		handlePtr(n)
	case reflect.Map:
		handleMap(n)
	}
}

func handlePtr(n *node) {
	ptrValue := reflect.ValueOf(n.data).Elem()

	if ptrValue.IsValid() {
		n.data = ptrValue.Interface()
		n.dataType = "&" + ptrValue.Type().Name()
		n.dataKind = ptrValue.Kind()
		handleValue(n)
	}
}

func handleStruct(n *node) {
	var wg sync.WaitGroup
	s := reflect.ValueOf(n.data)
	numberOfFields := s.NumField()

	for i := 0; i < numberOfFields; i++ {
		wg.Add(1)
		item := s.Field(i)
		itemName := s.Type().Field(i).Name

		for _, c := range itemName {
			if unicode.IsLower(c) {
				sValueCopy := reflect.New(s.Type()).Elem()
				sValueCopy.Set(s)
				itemTemp := sValueCopy.Field(i)
				item = reflect.NewAt(itemTemp.Type(), unsafe.Pointer(itemTemp.UnsafeAddr())).Elem()
			}

			break
		}

		newNode := n.append(&node{name: itemName, dataType: item.Type().Name(), dataKind: item.Kind(), data: item.Interface()})

		go func() {
			defer wg.Done()
			handleValue(newNode)
		}()
	}

	wg.Wait()
}

func handleSlice(n *node) {
	var wg sync.WaitGroup
	s := reflect.ValueOf(n.data)
	sLen := s.Len()
	n.dataType = s.Type().Elem().String()

	for i := 0; i < sLen; i++ {
		wg.Add(1)
		item := s.Index(i)
		newNode := n.append(&node{dataType: item.Type().Name(), dataKind: item.Kind(), data: item.Interface(), isSliceMember: true})

		go func() {
			defer wg.Done()
			handleValue(newNode)
		}()
	}

	wg.Wait()
}

func handleMap(n *node) {
	m := reflect.ValueOf(n.data)
	keys := m.MapKeys()
	n.dataType = fmt.Sprintf("map[%s]%s", m.Type().Key().String(), m.Type().Elem().String())

	for _, key := range keys {
		keyNode := &node{
			dataType: key.Type().Name(),
			dataKind: key.Kind(),
			data:     key.Interface(),
			isMapKey: true,
		}

		handleValue(keyNode)

		value := m.MapIndex(key)

		valueNode := &node{
			dataType:   value.Type().Name(),
			dataKind:   value.Kind(),
			data:       value.Interface(),
			isMapValue: true,
		}

		handleValue(valueNode)
		keyNode.associatedNode = valueNode
		n.append(keyNode)
	}
}
