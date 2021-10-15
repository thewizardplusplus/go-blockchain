package storing

import (
	"fmt"

	"github.com/thewizardplusplus/go-blockchain"
)

//go:generate mockery --name=Stringer --inpackage --case=underscore --testonly

// Stringer ...
//
// It's used only for mock generating.
//
type Stringer interface {
	fmt.Stringer
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
