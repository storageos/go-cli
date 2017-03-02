package opts

import "fmt"

// ValidateOperator validates whether the input string is a valid operator.
func ValidateOperator(val string) (string, error) {
	switch val {
	case "!", "=", "==", "in", "!=", "notin", "exists", "gt", "lt":
		return val, nil
	default:
		return "", fmt.Errorf("invalid operator %q.  Must be one of %q, %q, %q, %q, %q, %q, %q, %q, or %q", val, "!", "=", "==", "in", "!=", "notin", "exists", "gt", "lt")
	}
}

// ValidateRuleAction validates whether the input string is a valid action.
func ValidateRuleAction(val string) (string, error) {
	switch val {
	case "add", "remove":
		return val, nil
	default:
		return "", fmt.Errorf("invalid action: %s.  Must be one of %q or %q", val, "add", "remove")
	}
}
