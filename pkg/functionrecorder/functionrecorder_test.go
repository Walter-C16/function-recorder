package functionrecorder

import "testing"

type B struct {
	id          int
	AccountCode string
}

type A struct {
	ID                   string
	TestInnerStruct      B
	TestInnerEmptyStruct B
	TestInnerStructPtr   *B
	TestSlice            []string
	TestNilSlice         []string
}

func TestFunctionRecorder_Record(t *testing.T) {
	tests := []struct {
		name     string
		function interface{}
		args     []interface{}
	}{
		{
			name:     "Test",
			function: func(t A) A { return t },
			args: []interface{}{
				A{
					ID: "hello",
					TestInnerStruct: B{
						id:          1234,
						AccountCode: "codeRed",
					},
					TestInnerStructPtr: &B{
						id:          1234,
						AccountCode: "codeRed",
					},
					TestSlice: []string{"1", "hola mundo"},
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
