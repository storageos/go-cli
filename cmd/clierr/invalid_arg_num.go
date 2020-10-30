package clierr

import "fmt"

// ErrInvalidArgNum is an error raised when the command hasn't got the right
// amount of arguments. It contains info to explain the error and suggest the
// right way.
type ErrInvalidArgNum struct {
	got      []string
	expected int
	example  string
}

// Error returns a string representation of the error.
func (e *ErrInvalidArgNum) Error() string {
	return fmt.Sprintf("invalid number of arguments, got %v, expected %d.\nExample:\n  %s", e.got, e.expected, e.example)
}

// NewErrInvalidArgNum returns a new error when the command hasn't got the right
// amount of arguments. It contains info to explain the error and suggest the
// right way.
func NewErrInvalidArgNum(got []string, expected int, example string) *ErrInvalidArgNum {
	return &ErrInvalidArgNum{
		got:      got,
		expected: expected,
		example:  example,
	}
}
