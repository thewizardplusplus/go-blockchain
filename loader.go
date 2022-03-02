package blockchain

import (
	"github.com/pkg/errors"
)

// ErrNoMatch ...
var ErrNoMatch = errors.New("no match")

//go:generate mockery --name=Loader --inpackage --case=underscore --testonly

// Loader ...
type Loader interface {
	LoadBlocks(cursor interface{}, count int) (
		blocks BlockGroup,
		nextCursor interface{},
		err error,
	)
}

// FindDifferences ...
func FindDifferences(leftLoader Loader, rightLoader Loader, chunkSize int) (
	leftDifferences BlockGroup,
	rightDifferences BlockGroup,
	err error,
) {
	leftBlocks, _, err := leftLoader.LoadBlocks(nil, chunkSize)
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to load the left blocks")
	}

	rightBlocks, _, err := rightLoader.LoadBlocks(nil, chunkSize)
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to load the right blocks")
	}

	leftIndex, rightIndex, hasMatch := leftBlocks.FindDifferences(rightBlocks)
	if !hasMatch {
		return nil, nil, ErrNoMatch
	}

	return leftBlocks[:leftIndex], rightBlocks[:rightIndex], nil
}
