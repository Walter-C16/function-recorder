package functionrecorder

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unicode"
	"unsafe"
)

type FunctionRecorder struct {
	recordArguments, recordReturnedValues bool
}

const (
	_stringPattern  = "%s%s: \"%s\","
	_numericPattern = "%s%s: %s,"
)

const (
	_numeric = iota
	_string
)

func NewFunctionRecorder(recordArguments, recordReturnedValues bool) FunctionRecorder {
	return FunctionRecorder{recordArguments: recordArguments, recordReturnedValues: recordReturnedValues}
}

func (fr FunctionRecorder) Record(function interface{}, args ...interface{}) {

	a := make([]reflect.Value, len(args))

	for i := range args {
		a[i] = reflect.ValueOf(args[i])
	}

	returns := reflect.ValueOf(function).Call(a)

	if fr.recordArguments {
		fmt.Println("Arguments:")
		argsLength := len(args)

		for i := range args {
			value := reflect.ValueOf(args[i])
			handleValue(args[i], value.Kind().String(), value.Type().Name(), value.Kind(), 0, i == argsLength-1)
		}
	}

	if fr.recordReturnedValues {
		fmt.Println("Returned values:")
		returnsLength := len(returns)

		for i := range returns {
			value := reflect.ValueOf(args[i])
			handleValue(args[i], value.Kind().String(), value.Type().Name(), value.Kind(), 0, i == returnsLength-1)
		}
	}
}

func handleValue(value interface{}, valueName, valueType string, kind reflect.Kind, level int, isSliceMember bool) {
	switch kind {
	case reflect.Struct:
		handleStruct(value, valueName, valueType, level)
	case reflect.String:
		print(value.(string), valueName, _string, level, isSliceMember)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64:
		print(strconv.Itoa(value.(int)), valueName, _numeric, level, isSliceMember)
	case reflect.Ptr:
		handlePtr(value, valueName, valueType, level)
	case reflect.Array, reflect.Slice:
		handleSlice(value, valueName, level)
	}
}

func handleSlice(slice interface{}, sliceName string, level int) {
	items := reflect.ValueOf(slice)
	tabbedString := getTabbedString(level)
	sliceLength := items.Len()

	if sliceLength == 0 {
		fmt.Printf("%s%s: []%s{},\n", tabbedString, sliceName, reflect.TypeOf(slice).Elem().String())
		return
	}

	fmt.Printf("%s%s: []%s{\n", tabbedString, sliceName, reflect.TypeOf(slice).Elem().String())

	for i := 0; i < sliceLength; i++ {
		item := items.Index(i)
		handleValue(item.Interface(), "", "", item.Kind(), level+1, true)
	}

	fmt.Printf("%s},\n", tabbedString)
}

func handlePtr(ptr interface{}, ptrName, ptrType string, level int) {
	ptrValue := reflect.ValueOf(ptr).Elem()

	if ptrValue.IsValid() {
		handleValue(ptrValue.Interface(), ptrName, "&"+ptrValue.Type().Name(), ptrValue.Kind(), level, false)
	}
}

func handleStruct(st interface{}, structName, structType string, level int) {
	structValue := reflect.ValueOf(st)
	numField := structValue.NumField()
	tabbedString := getTabbedString(level)

	if numField == 0 {
		fmt.Printf("%s%s: %s {},\n", tabbedString, structName, structType)
		return
	}

	fmt.Printf("%s%s: %s {\n", tabbedString, structName, structType)

	for i := 0; i < numField; i++ {
		field := structValue.Field(i)
		fieldName := structValue.Type().Field(i).Name

		for _, c := range fieldName {
			if unicode.IsLower(c) {
				structValueCopy := reflect.New(structValue.Type()).Elem()
				structValueCopy.Set(structValue)
				fieldTemp := structValueCopy.Field(i)
				field = reflect.NewAt(fieldTemp.Type(), unsafe.Pointer(fieldTemp.UnsafeAddr())).Elem()
			}

			break
		}

		handleValue(field.Interface(), fieldName, field.Type().Name(), field.Kind(), level+1, false)
	}

	fmt.Printf("%s},\n", tabbedString)
}

func print(value, valueName string, valueType, level int, isSliceMember bool) {
	var strToPrint, pattern string
	tabbedString := getTabbedString(level)

	switch valueType {
	case _numeric:
		pattern = _numericPattern
	case _string:
		pattern = _stringPattern
	}

	strToPrint = fmt.Sprintf(pattern, tabbedString, valueName, value)

	if isSliceMember {
		strToPrint = strings.Replace(strToPrint, ": ", "", 1)
	}

	fmt.Println(strToPrint)
}

func getTabbedString(n int) (tabbedString string) {
	for i := 0; i < n; i++ {
		tabbedString = tabbedString + "\t"
	}

	return
}
