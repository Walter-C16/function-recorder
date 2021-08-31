package functionrecorder

import (
	"testing"
)

type B struct {
	id          int
	AccountCode string
}

type A struct {
	ID                   string
	TestInnerStruct      B
	TestInnerEmptyStruct B
	TestInnerPtrStruct   *B
	TestSlice            []string
	TestNilSlice         []string
	TestArray            [3]int
	TestMap              map[int]string
	TestNilMap           map[string]int
	TestInt8             int8
	TestInt16            int16
	TestInt32            int32
	TestInt64            int64
	TestUint             uint
	TestUint8            uint8
	TestUInt16           uint16
	TestUInt32           uint32
	TestUInt64           uint64
	TestBool             bool
	TestFloat32          float32
	TestFloat64          float64
	TestComplex64        complex64
	TestComplex128       complex128
}

func TestFunctionRecorder_Record(t *testing.T) {
	tests := []struct {
		name     string
		function interface{}
		args     []interface{}
	}{
		{
			name:     "Test",
			function: func(t A, c *B) A { return t },
			args: []interface{}{
				A{
					ID: "7d4353e7-5194-4f54-9c45-994a5caf42ec",
					TestInnerStruct: B{
						id:          1234,
						AccountCode: "codeRed",
					},
					TestInnerPtrStruct: &B{
						id:          1234,
						AccountCode: "codeRed",
					},
					TestSlice: []string{"1", "hola mundo"},
					TestArray: [3]int{1, 2, 4},
					TestMap:   map[int]string{1: "hola", 3: "mundo"},
				},
				&B{
					id:          12345,
					AccountCode: "codeBlue",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fr := NewFunctionRecorder(true, true)
			fr.Record(tt.function, tt.args...)
		})
	}
}
