package service

import (
	"capi/domain"
	"capi/dto"
	"capi/errs"
)

type CustomerService interface {
	GetAllCustomers(string) ([]dto.CustomerResponse, *errs.AppErr)
	GetCustomerByID(string) (*dto.CustomerResponse, *errs.AppErr)
}

type DefaultCustomerService struct {
	repo domain.CustomerRepository
}

func (s DefaultCustomerService) GetAllCustomers(status string) ([]dto.CustomerResponse, *errs.AppErr) {
	if status == "active" {
		status = "1"
	} else if status == "inactive" {
		status = "0"
	} else if status == "" {
		status = ""
	} else {
		return nil, errs.NewBadRequestError("invalid request")
	}

	customers, err := s.repo.FindAll(status)
	if err != nil {
		return nil, err
	}

	var response []dto.CustomerResponse
	for _, customer := range customers {
		response = append(response, customer.ToDTO())
	}

	return response, nil
}

func (s DefaultCustomerService) GetCustomerByID(CustomerID string) (*dto.CustomerResponse, *errs.AppErr) {
	cust, err := s.repo.FindByID(CustomerID)
	if err != nil {
		return nil, err
	}

	response := cust.ToDTO()

	return &response, nil
}

func NewCustomerService(repository domain.CustomerRepository) DefaultCustomerService {
	return DefaultCustomerService{repository}
}