package functionrecorder

import "testing"

type TestInnerStruct struct {
	id          int
	AccountCode string
}

type TestStruct struct {
	ID                   string
	PersonalInformation  TestInnerStruct
	PersonalInformation2 TestInnerStruct
}

func TestFunctionRecorder_Record(t *testing.T) {
	tests := []struct {
		name     string
		function interface{}
		args     []interface{}
	}{
		{
			name:     "Test",
			function: func(t TestStruct) TestStruct { return t },
			args: []interface{}{
				TestStruct{
					ID: "hello",
					PersonalInformation: TestInnerStruct{
						id:          1234,
						AccountCode: "codeRed",
					},
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
