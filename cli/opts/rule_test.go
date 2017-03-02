package opts

import (
	"fmt"
	"testing"
)

func TestValidateOperator(t *testing.T) {

	e := fmt.Errorf("some error")
	var tests = []struct {
		input       string
		expected    string
		expectedErr error
	}{
		{"!", "!", nil},
		{"=", "=", nil},
		{"==", "==", nil},
		{"in", "in", nil},
		{"!=", "!=", nil},
		{"notin", "notin", nil},
		{"exists", "exists", nil},
		{"gt", "gt", nil},
		{"lt", "lt", nil},
		{"foo", "", e},
		{"===", "", e},
		{" ", "", e},
		{"  ", "", e},
	}
	for _, tt := range tests {
		output, err := ValidateOperator(tt.input)
		if err != nil && tt.expectedErr == nil {
			t.Errorf("ValidateOperator(%q): Got unexpected error %q.", tt.input, err)
		}
		if output != tt.expected {
			t.Errorf("ValidateOperator(%q): Got %q. Want %q.", tt.input, output, tt.expected)
		}
	}
}

func TestRuleAction(t *testing.T) {

	e := fmt.Errorf("some error")
	var tests = []struct {
		input       string
		expected    string
		expectedErr error
	}{
		{"add", "add", nil},
		{"remove", "remove", nil},
		{"foo", "", e},
		{"in", "", e},
		{"!=", "", e},
		{"notin", "", e},
		{"exists", "", e},
		{" ", "", e},
		{"  ", "", e},
	}
	for _, tt := range tests {
		output, err := ValidateRuleAction(tt.input)
		if err != nil && tt.expectedErr == nil {
			t.Errorf("ValidateRuleAction(%q): Got unexpected error %q.", tt.input, err)
		}
		if output != tt.expected {
			t.Errorf("ValidateRuleAction(%q): Got %q. Want %q.", tt.input, output, tt.expected)
		}
	}
}
