package blockchain

import (
	"fmt"
	"reflect"
)

//go:generate mockery --name=DataComparer --inpackage --case=underscore --testonly

// DataComparer ...
type DataComparer interface {
	Equal(data Data) bool
}

//go:generate mockery --name=Data --inpackage --case=underscore --testonly

// Data ...
type Data interface {
	fmt.Stringer
	DataComparer
}

type universalDataWrapper struct {
	innerData interface{}
}

func (wrapper universalDataWrapper) String() string {
	// explicitly check this interface to prioritize its use
	if stringer, ok := wrapper.innerData.(fmt.Stringer); ok {
		return stringer.String()
	}

	return fmt.Sprint(wrapper.innerData)
}

func (wrapper universalDataWrapper) MarshalText() ([]byte, error) {
	text := wrapper.String()
	return []byte(text), nil
}

func (wrapper universalDataWrapper) Equal(data Data) bool {
	if comparer, ok := wrapper.innerData.(DataComparer); ok {
		return comparer.Equal(data)
	}

	if anotherWrapper, ok := data.(universalDataWrapper); ok {
		return reflect.DeepEqual(wrapper.innerData, anotherWrapper.innerData)
	}

	return false
}
