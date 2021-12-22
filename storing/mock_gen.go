package storing

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

//go:generate mockery --name=Storage --inpackage --case=underscore --testonly

// Storage ...
//
// It's used only for mock generating.
//
type Storage interface {
	blockchain.Storage
}

//go:generate mockery --name=GroupStorage --inpackage --case=underscore --testonly

// GroupStorage ...
//
// It's used only for mock generating.
//
type GroupStorage interface {
	blockchain.GroupStorage
}
