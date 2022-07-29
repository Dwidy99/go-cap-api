package domain

import (
	"capi/dto"
	"capi/errs"
)

type Customer struct {
	ID          string `json:"id" xml:"id" db:"customer_id"`
	Name        string `json:"name" xml:"name"`
	City        string `json:"city" xml:"city"`
	ZipCode     string `json:"zip_code" xml:"zipCode"`
	DateOfBirth string `json:"date_of_birth" xml:"date_of_birth" db:"date_of_birth"`
	Status      string `json:"status" xml:"status"`
}

// Kontrak untuk di implementasikan di package service
type CustomerRepository interface {
	FindAll(string) ([]Customer, *errs.AppErr)
	FindByID(string) (*Customer, *errs.AppErr)
}

func (c Customer) convertStatusCustomer() string {
	statusCustomer := "active"
	if c.Status == "0" {
		statusCustomer = "inactive"
	}

	return statusCustomer
}

// dto.CustomerResponse ada di package dto
func (c Customer) ToDTO() dto.CustomerResponse {

	return dto.CustomerResponse{
		ID: c.ID,
		Name: c.Name,
		DateOfBirth: c.DateOfBirth,
		City: c.City,
		ZipCode: c.ZipCode,
		Status: c.convertStatusCustomer(),
	}
}
