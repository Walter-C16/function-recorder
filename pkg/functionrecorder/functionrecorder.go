package functionrecorder

import (
	"fmt"
	"reflect"
	"strings"
)

type FunctionRecorder struct {
	recordArguments, recordReturnedValues bool
}

const (
	defaultPattern     = "%s%s: %v,\n"
	stringPattern      = "%s%s: \"%s\",\n"
	arrayPattern       = "%s%s: [%d]%s{\n"
	emptyArrayPattern  = "%s%s: [0]%s{},\n"
	slicePattern       = "%s%s: []%s{\n"
	emptySlicePattern  = "%s%s: []%s{},\n"
	structPattern      = "%s%s: %s{\n"
	emptyStructPattern = "%s%s: %s{},\n"
	mapPattern         = "%s%s: %s{\n"
	emptyMapPattern    = "%s%s: %s{},\n"
)

func NewFunctionRecorder(recordArguments, recordReturnedValues bool) FunctionRecorder {
	return FunctionRecorder{recordArguments: recordArguments, recordReturnedValues: recordReturnedValues}
}

func (fr FunctionRecorder) Record(function interface{}, args ...interface{}) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Print("Recovered: ")
			fmt.Println(r)
		}
	}()

	a := make([]reflect.Value, len(args))

	for i := range args {
		a[i] = reflect.ValueOf(args[i])
	}

	returns := reflect.ValueOf(function).Call(a)

	if fr.recordArguments {
		fmt.Println("Arguments:")
		printTrees(getTrees(args))
	}

	if fr.recordReturnedValues {
		fmt.Println("Returns:")
		r := make([]interface{}, len(returns))

		for i := range returns {
			r[i] = returns[i].Interface()
		}

		printTrees(getTrees(r))
	}
}

func getTrees(values []interface{}) []*valueTree {
	numberOfValues := len(values)
	argumentTrees := make([]*valueTree, 0, numberOfValues)
	for i := range values {
		argumentTrees = append(argumentTrees, buildTree(values[i]))
	}

	return argumentTrees
}

func printTrees(trees []*valueTree) {
	var output strings.Builder
	for i := range trees {
		if trees[i].root == nil {
			continue
		}

		trees[i].root.name = trees[i].root.dataKind.String()
		printNode(&output, trees[i].root, 0)
	}

	fmt.Println(output.String())
}

func printNode(s *strings.Builder, n *node, level int) {
	switch n.dataKind {
	case reflect.String, reflect.Int, reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128, reflect.Int16, reflect.Int32,
		reflect.Int64, reflect.Int8, reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint8, reflect.Bool:
		printPrimitive(s, n, level)
	case reflect.Array, reflect.Slice, reflect.Struct, reflect.Map:
		printComposite(s, n, level)
	default:
		for i := range n.children {
			printNode(s, n.children[i], level+1)
		}
	}

	if n.associatedNode != nil {
		printNode(s, n.associatedNode, level)
	}
}

func printComposite(s *strings.Builder, n *node, level int) {
	var emptyPattern, pattern string
	tabbedString := getTabbedString(level)

	switch n.dataKind {
	case reflect.Array:
		pattern = arrayPattern
		emptyPattern = emptyArrayPattern
	case reflect.Slice:
		pattern = slicePattern
		emptyPattern = emptySlicePattern
	case reflect.Struct:
		pattern = structPattern
		emptyPattern = emptyStructPattern
	case reflect.Map:
		pattern = mapPattern
		emptyPattern = emptyMapPattern
	}

	if len(n.children) == 0 {
		s.WriteString(fmt.Sprintf(emptyPattern, tabbedString, n.name, n.dataType))
		return
	}

	if n.dataKind == reflect.Array {
		s.WriteString(fmt.Sprintf(pattern, tabbedString, n.name, len(n.children), n.dataType))
	} else {
		s.WriteString(fmt.Sprintf(pattern, tabbedString, n.name, n.dataType))
	}

	for i := range n.children {
		printNode(s, n.children[i], level+1)
	}

	s.WriteString(fmt.Sprintf("%s},\n", tabbedString))
}

func printPrimitive(s *strings.Builder, n *node, level int) {
	var str string

	pattern := defaultPattern

	if n.dataKind == reflect.String {
		pattern = stringPattern
	}

	str = fmt.Sprintf(pattern, getTabbedString(level), n.name, n.data)

	switch {
	case n.isMapKey:
		str = strings.Replace(strings.Replace(str, ": ", "", 1), ",\n", " : ", 1)
	case n.isMapValue:
		str = strings.Replace(strings.Replace(str, ": ", "", 1), getTabbedString(level), "", 1)
	case n.isSliceMember:
		str = strings.Replace(str, ": ", "", 1)
	}

	s.WriteString(str)
}

func getTabbedString(level int) (tabbedString string) {
	for i := 0; i < level; i++ {
		tabbedString += "\t"
	}

	return
}
