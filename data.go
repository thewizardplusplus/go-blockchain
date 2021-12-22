package blockchain

import (
	"fmt"
)

//go:generate mockery --name=Data --inpackage --case=underscore --testonly

// Data ...
type Data interface {
	fmt.Stringer

	Equal(data Data) bool
}
