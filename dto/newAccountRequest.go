package dto

import (
	"capi/errs"
	"strings"
)

const (
	ErrAccountRequestAmount      = "minimum amount to open new account is 5000"
	ErrAccountRequestAccountType = "account type must be checking or saving"
)

type NewAccountRequest struct {
	CustomerID  string  `json:"customer_id"`
	AccountType string  `json:"account_type"`
	Amount      float64 `json:"amount"`
}

func (r NewAccountRequest) Validate() *errs.AppErr {
	if r.Amount < 5000 {
		// return errs.NewValidationError("minimum amount to open new account is 5000")
		return errs.NewValidationError(ErrAccountRequestAmount)
	}
	if strings.ToLower(r.AccountType) != "saving" && strings.ToLower(r.AccountType) != "checking" {
		// return errs.NewValidationError("account type must be checking or saving")
		return errs.NewValidationError(ErrAccountRequestAccountType)
	}

	return nil
}