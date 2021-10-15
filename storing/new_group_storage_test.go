package storing

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/go-blockchain"
)

func TestNewGroupStorage(test *testing.T) {
	type args struct {
		storage blockchain.Storage
	}

	for _, data := range []struct {
		name             string
		args             args
		wantGroupStorage blockchain.GroupStorage
	}{
		{
			name: "with the group storage",
			args: args{
				storage: new(MockGroupStorage),
			},
			wantGroupStorage: new(MockGroupStorage),
		},
		{
			name: "with the storage",
			args: args{
				storage: new(MockStorage),
			},
			wantGroupStorage: GroupStorageWrapper{Storage: new(MockStorage)},
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			gotGroupStorage := NewGroupStorage(data.args.storage)

			mock.AssertExpectationsForObjects(test, data.args.storage)
			assert.Equal(test, data.wantGroupStorage, gotGroupStorage)
		})
	}
}
