package dto

import (
	"capi/errs"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateTrasaction_valid_transaction(t *testing.T) {
	tx := TransactionRequest{
		Amount:          1,
		TransactionType: WITHDRAWAL,
	}

	res := tx.Validate()

	assert.Nil(t, res)
}

func TestValidateTrasaction_invalid_transaction_type(t *testing.T) {
	tx := TransactionRequest{
		Amount:          1,
		TransactionType: "invalid type",
	}

	res := tx.Validate()

	expect := errs.NewValidationError("")

	assert.NotNil(t, res)
	assert.Equal(t, expect.Code, res.Code)
}

func TestValidateTrasaction_invalid_transaction_amount(t *testing.T) {
	tx := TransactionRequest{
		Amount:          -1,
		TransactionType: WITHDRAWAL,
	}

	res := tx.Validate()

	expect := errs.NewValidationError("")

	assert.NotNil(t, res)
	assert.Equal(t, expect.Code, res.Code)
}