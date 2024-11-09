package armparser

import "fmt"

type ArgumentError struct {
	function string
	expected int
	got      int
}

func (e *ArgumentError) Error() string {
	return fmt.Sprintf("%s function expects %d arguments but got %d", e.function, e.expected, e.got)
}

func NewArgumentError(function string, expected, got int) error {
	return &ArgumentError{function: function, expected: expected, got: got}
}
