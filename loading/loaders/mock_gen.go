package loaders

import (
	"github.com/thewizardplusplus/go-blockchain"
)

//go:generate mockery --name=Data --inpackage --case=underscore --testonly

// Data ...
//
// It's used only for mock generating.
//
type Data interface {
	blockchain.Data
}
