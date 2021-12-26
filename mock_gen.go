package blockchain

import (
	"encoding"
	"fmt"
)

//go:generate mockery --name=Stringer --inpackage --case=underscore --testonly

// Stringer ...
//
// It's used only for mock generating.
//
type Stringer interface {
	fmt.Stringer
}

//go:generate mockery --name=TextMarshaler --inpackage --case=underscore --testonly

// TextMarshaler ...
//
// It's used only for mock generating.
//
type TextMarshaler interface {
	encoding.TextMarshaler
}
