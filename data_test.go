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

			// nolint: staticcheck
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

			// nolint: staticcheck
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

func Test_universalDataWrapper_Equal(test *testing.T) {
	type fields struct {
		innerData interface{}
	}
	type args struct {
		data Data
	}

	for _, data := range []struct {
		name   string
		fields fields
		args   args
		want   assert.BoolAssertionFunc
	}{
		{
			name: "inner data is DataComparer/false",
			fields: fields{
				innerData: func() interface{} {
					stringer := new(MockDataComparer)
					stringer.
						On("Equal", mock.AnythingOfType("*blockchain.MockData")).
						Return(false)

					return stringer
				}(),
			},
			args: args{
				data: new(MockData),
			},
			want: assert.False,
		},
		{
			name: "inner data is DataComparer/true",
			fields: fields{
				innerData: func() interface{} {
					stringer := new(MockDataComparer)
					stringer.
						On("Equal", mock.AnythingOfType("*blockchain.MockData")).
						Return(true)

					return stringer
				}(),
			},
			args: args{
				data: new(MockData),
			},
			want: assert.True,
		},
		{
			name: "data is universalDataWrapper/false",
			fields: fields{
				innerData: "test #1",
			},
			args: args{
				data: universalDataWrapper{
					innerData: "test #2",
				},
			},
			want: assert.False,
		},
		{
			name: "data is universalDataWrapper/true",
			fields: fields{
				innerData: "test",
			},
			args: args{
				data: universalDataWrapper{
					innerData: "test",
				},
			},
			want: assert.True,
		},
		{
			name: "inner data is not DataComparer and data is not universalDataWrapper",
			fields: fields{
				innerData: "test",
			},
			args: args{
				data: new(MockData),
			},
			want: assert.False,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			wrapper := universalDataWrapper{
				innerData: data.fields.innerData,
			}
			got := wrapper.Equal(data.args.data)

			for _, mockInstance := range []interface{}{
				data.fields.innerData,
				data.args.data,
			} {
				// nolint: staticcheck
				if _, ok := mockInstance.(interface {
					AssertExpectations(mock.TestingT) bool
				}); ok {
					mock.AssertExpectationsForObjects(test, mockInstance)
				}
			}
			data.want(test, got)
		})
	}
}
