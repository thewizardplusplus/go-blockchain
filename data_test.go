package blockchain

import (
	"testing"
	"testing/iotest"

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

func Test_universalDataWrapper_MarshalText(test *testing.T) {
	type fields struct {
		innerData interface{}
	}

	for _, data := range []struct {
		name     string
		fields   fields
		wantText []byte
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			name: "integer",
			fields: fields{
				innerData: 23,
			},
			wantText: []byte("23"),
			wantErr:  assert.NoError,
		},
		{
			name: "string",
			fields: fields{
				innerData: "test",
			},
			wantText: []byte("test"),
			wantErr:  assert.NoError,
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
			wantText: []byte("test"),
			wantErr:  assert.NoError,
		},
		{
			name: "encoding.TextMarshaler/success",
			fields: fields{
				innerData: func() interface{} {
					stringer := new(MockTextMarshaler)
					stringer.On("MarshalText").Return([]byte("test"), nil)

					return stringer
				}(),
			},
			wantText: []byte("test"),
			wantErr:  assert.NoError,
		},
		{
			name: "encoding.TextMarshaler/error",
			fields: fields{
				innerData: func() interface{} {
					stringer := new(MockTextMarshaler)
					stringer.On("MarshalText").Return(nil, iotest.ErrTimeout)

					return stringer
				}(),
			},
			wantText: nil,
			wantErr:  assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			wrapper := universalDataWrapper{
				innerData: data.fields.innerData,
			}
			gotText, gotErr := wrapper.MarshalText()

			if _, ok := data.fields.innerData.(interface {
				AssertExpectations(mock.TestingT) bool
			}); ok {
				mock.AssertExpectationsForObjects(test, data.fields.innerData)
			}
			assert.Equal(test, data.wantText, gotText)
			data.wantErr(test, gotErr)
		})
	}
}
