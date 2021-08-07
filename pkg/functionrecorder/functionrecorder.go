package functionrecorder

import (
	"fmt"
	"reflect"
	"strconv"
	"unicode"
	"unsafe"
)

type FunctionRecorder struct {
	recordArguments, recordReturnedValues bool
}

const (
	Invalid int = iota
	Bool
	Numeric
	Complex64
	Complex128
	Array
	Chan
	Func
	Interface
	Map
	Ptr
	Slice
	String
	Struct
	UnsafePointer
)

const (
	_stringPattern = ""
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
			handleValue(args[i], "", value.Kind(), 0, i == argsLength-1)
		}
	}

	if fr.recordReturnedValues {
		fmt.Println("Returned values:")
		returnsLength := len(returns)

		for i := range returns {
			value := reflect.ValueOf(args[i])
			handleValue(args[i], "", value.Kind(), 0, i == returnsLength-1)
		}
	}
}

func handleValue(v interface{}, name string, kind reflect.Kind, level int, isLastValue bool) {
	switch kind {
	case reflect.Struct:
		vValue := reflect.ValueOf(v)
		handleStruct(v, vValue.Kind().String(), vValue.Type().Name(), level, isLastValue)
	case reflect.String:
		printValues(reflect.ValueOf(v).Interface().(string), "", String, level, isLastValue)
	}
}

func handleStruct(st interface{}, fieldName, structName string, level int, isStructLastFieldOrLastValue bool) {
	v := reflect.ValueOf(st)
	numField := v.NumField()
	tabbedString := getTabbedString(level)

	if fieldName != "" {
		fmt.Printf("%s%s: %s {\n", tabbedString, fieldName, structName)
	} else {
		fmt.Printf("%s{\n", tabbedString)
	}

	for i := 0; i < numField; i++ {
		field := v.Field(i)
		fieldName := v.Type().Field(i).Name

		for _, c := range fieldName {
			if unicode.IsLower(c) {
				vCopy := reflect.New(v.Type()).Elem()
				vCopy.Set(v)
				fieldTemp := vCopy.Field(i)
				field = reflect.NewAt(fieldTemp.Type(), unsafe.Pointer(fieldTemp.UnsafeAddr())).Elem()
			}

			break
		}

		switch field.Kind() {
		case reflect.Struct:
			handleStruct(field.Interface(), fieldName, field.Type().Name(), level+1, i == numField-1)
		case reflect.String:
			printValues(field.Interface().(string), fieldName, String, level+1, i == numField-1)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64:
			printValues(strconv.Itoa(field.Interface().(int)), fieldName, Numeric, level+1, i == numField-1)
		}
	}

	stTrailingToPrint := fmt.Sprintf("%s}", tabbedString)

	if !isStructLastFieldOrLastValue {
		stTrailingToPrint = stTrailingToPrint + ","
	}

	fmt.Println(stTrailingToPrint)
}

func printValues(value, fieldName string, valueType, level int, isStructLastFieldOrLastValue bool) {
	var strToPrint string
	tabbedString := getTabbedString(level)

	switch valueType {
	case String:
		if fieldName != "" {
			strToPrint = fmt.Sprintf("%s%s: \"%s\"", tabbedString, fieldName, value)
		} else {
			strToPrint = fmt.Sprintf("%s\"%s\"", tabbedString, value)
		}
	case Numeric:
		if fieldName != "" {
			strToPrint = fmt.Sprintf("%s%s: %s", tabbedString, fieldName, value)
		} else {
			strToPrint = fmt.Sprintf("%s%s", tabbedString, value)
		}
	}

	if !isStructLastFieldOrLastValue {
		strToPrint = strToPrint + ","
	}

	fmt.Println(strToPrint)
}

func getTabbedString(n int) (tabbedString string) {
	for i := 0; i < n; i++ {
		tabbedString = tabbedString + "\t"
	}

	return
}
