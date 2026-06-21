package validate

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testRequest struct {
	Name string `validate:"required,min=3"`
	Age  int    `validate:"gte=18"`
}

func TestNewValidator(t *testing.T) {
	v, err := NewValidator()

	require.NoError(t, err)
	require.NotNil(t, v)
}

func TestValidate(t *testing.T) {
	v, err := NewValidator()
	require.NoError(t, err)

	tests := []struct {
		name      string
		req       testRequest
		wantError bool
	}{
		{
			name: "valid",
			req: testRequest{
				Name: "Alice",
				Age:  18,
			},
			wantError: false,
		},
		{
			name: "required",
			req: testRequest{
				Name: "",
				Age:  18,
			},
			wantError: true,
		},
		{
			name: "min",
			req: testRequest{
				Name: "ab",
				Age:  18,
			},
			wantError: true,
		},
		{
			name: "gte",
			req: testRequest{
				Name: "Alice",
				Age:  17,
			},
			wantError: true,
		},
		{
			name: "multiple errors",
			req: testRequest{
				Name: "",
				Age:  10,
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Validate(tt.req)

			if tt.wantError {
				require.Error(t, err)

				_, ok := err.(validator.ValidationErrors)
				assert.True(t, ok)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
