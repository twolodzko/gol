package parser

import (
	"fmt"
)

// Symbol is a generic type for a named object
type Symbol struct {
	Name string
}

// Print symbols unquoted, contrary to strings
func (s Symbol) String() string {
	return s.Name
}

// String is a custom string type
type String struct {
	Characters string
}

// Print strings quoted
func (s String) String() string {
	return fmt.Sprintf("\"%s\"", s.Characters)
}
