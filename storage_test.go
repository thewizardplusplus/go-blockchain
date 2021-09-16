package blockchain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLoadStorage(test *testing.T) {
	type args struct {
		storage       Storage
		loader        Loader
		initialCursor interface{}
		chunkSize     int
	}

	for _, data := range []struct {
		name           string
		args           args
		wantLastCursor interface{}
		wantErr        assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	} {
		test.Run(data.name, func(test *testing.T) {
			gotLastCursor, gotErr := LoadStorage(
				data.args.storage,
				data.args.loader,
				data.args.initialCursor,
				data.args.chunkSize,
			)

			mock.AssertExpectationsForObjects(test, data.args.storage, data.args.loader)
			assert.Equal(test, data.wantLastCursor, gotLastCursor)
			data.wantErr(test, gotErr)
		})
	}
}
