package loaders

import (
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
