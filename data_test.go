package blockchain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewData(test *testing.T) {
	data := NewData("test")

	wantedData := universalDataWrapper{innerData: "test"}
	assert.Equal(test, wantedData, data)
}

func Test_universalDataWrapper_String(test *testing.T) {
	type fields struct {
		innerData interface{}
	}

	for _, data := range []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "integer",
			fields: fields{
				innerData: 23,
			},
			want: "23",
		},
		{
			name: "string",
			fields: fields{
				innerData: "test",
			},
			want: "test",
		},
		{
			name: "fmt.Stringer",
			fields: fields{
				innerData: func() interface{} {
					stringer := new(MockStringer)
					stringer.On("String").Return("test")

					return stringer
				}(),
			},
			want: "test",
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			wrapper := universalDataWrapper{
				innerData: data.fields.innerData,
			}
			got := wrapper.String()

			if _, ok := data.fields.innerData.(interface {
				AssertExpectations(mock.TestingT) bool
			}); ok {
				mock.AssertExpectationsForObjects(test, data.fields.innerData)
			}
			assert.Equal(test, data.want, got)
		})
	}
}
