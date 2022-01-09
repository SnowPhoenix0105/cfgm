package obj2tree

import "fmt"

type PointerError struct {
	Level int
	Type  string
}

func newPointerError(level int, typ string) *PointerError {
	return &PointerError{
		Level: level,
		Type:  typ,
	}
}

func (err *PointerError) Error() string {
	return fmt.Sprintf("%d-level pointer is not allowed for %s.", err.Level, err.Type)
}
