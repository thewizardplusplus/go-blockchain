package blockchain

import (
	"fmt"
)

// Data ...
type Data interface {
	fmt.Stringer

	Equal(data Data) bool
}
