package proofers

import (
	"github.com/thewizardplusplus/go-blockchain"
)

//go:generate mockery --name=Hasher --inpackage --case=underscore --testonly

// Hasher ...
//
// It's used only for mock generating.
//
type Hasher interface {
	blockchain.Hasher
}
