package loading

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

//go:generate mockery --name=Proofer --inpackage --case=underscore --testonly

// Proofer ...
//
// It's used only for mock generating.
//
type Proofer interface {
	blockchain.Proofer
}

//go:generate mockery --name=Loader --inpackage --case=underscore --testonly

// Loader ...
//
// It's used only for mock generating.
//
type Loader interface {
	blockchain.Loader
}

//go:generate mockery --name=GroupStorage --inpackage --case=underscore --testonly

// GroupStorage ...
//
// It's used only for mock generating.
//
type GroupStorage interface {
	blockchain.GroupStorage
}
