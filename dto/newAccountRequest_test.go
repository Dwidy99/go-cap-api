package dto

import (
	"capi/errs"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccountRequestValidation(t *testing.T) {
	testCase := []struct {
		name   string
		input  NewAccountRequest
		expect *errs.AppErr
	}{
		{
			"valid request",
			NewAccountRequest{
				Amount:      5000,
				AccountType: "saving",
			},
			nil,
		},
		{
			"invalid amount",
			NewAccountRequest{
				Amount:      1,
				AccountType: "saving",
			},
			errs.NewValidationError(ErrAccountRequestAmount),
		},
		{
			"invalid account type",
			NewAccountRequest{
				Amount:      5000,
				AccountType: "invalid account type",
			},
			errs.NewValidationError(ErrAccountRequestAccountType),
		},
	}

	for _, tt := range testCase {
		t.Run(tt.name, func(t *testing.T) {
			res := tt.input.Validate()

			assert.Equal(t, res, tt.expect)
		})
	}
}