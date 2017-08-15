package opts

import (
	"fmt"
	"testing"
)

func TestValidateClusterSize(t *testing.T) {

	e := fmt.Errorf("some error")
	var tests = []struct {
		input       int
		expected    int
		expectedErr error
	}{
		{1, 1, nil},
		{3, 3, nil},
		{5, 5, nil},
		{7, 7, nil},
		{0, 0, e},
		{2, 0, e},
		{4, 0, e},
		{6, 0, e},
		{8, 0, e},
		{100, 0, e},
	}
	for _, tt := range tests {
		output, err := ValidateClusterSize(tt.input)
		if err != nil && tt.expectedErr == nil {
			t.Errorf("ValidateOperator(%d): Got unexpected error %d.", tt.input, err)
		}
		if output != tt.expected {
			t.Errorf("ValidateOperator(%d): Got %d. Want %d.", tt.input, output, tt.expected)
		}
	}
}
