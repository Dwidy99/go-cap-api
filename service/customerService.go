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

// Penjelasan DB 2.1
// struct DefaultCustomerService berisi properti repository ber-tipe **interface 
// CustomerRepository di dalam package domain
type DefaultCustomerService struct {
	repository domain.CustomerRepository
}

// penjelasan DB 2
// NewCustomerService(parameter) akan menerima value dari NewCustomerRepositoryDB(), value tersebut 
// akan mengisi properti repository yang ada di dalam struct DefaultCustomerService ber-tipe 
// domain.CustomerRepository -> (interface yang ada di package domain file customer.go)
func NewCustomerService(repository domain.CustomerRepository) DefaultCustomerService {
	return DefaultCustomerService{repository: repository}
}

// Penjelasan AllCustomers 2.
// Method GetAllCustomers(customerStatus string) mengembalikan semua data customers berupa array / slice yang 
// akan di gunakan di getAllCustomers() di file handle.go
func (s DefaultCustomerService) GetAllCustomers(customerStatus string) ([]dto.CustomerResponse, *errs.AppErr) {
	// FindAll(customerStatus) adalah kembalian semua data berdasarkan request di Browser berdasarkan argument customerStatus
	customers, err := s.repository.FindAll(customerStatus)
	// Cek Error
	if err != nil {
		return nil, errs.NewUnexpectedError("Unexpected DataBase Status")
	}
	var dtoCustomers []dto.CustomerResponse
	// Looping data customers dengan meng-convert field "status" dari int ke string
	for _, customer := range customers {
		dtoCustomers =append(dtoCustomers, customer.ToDTO())
	}

	// Kembalikan semua data customers berupa array/slice dan error
	return dtoCustomers, nil
}

// Penjelasan ByID 2.
// Method GetCustomerByID() mengembalikan data bukan json
func (s DefaultCustomerService) GetCustomerByID(CustomerID string) (*dto.CustomerResponse, *errs.AppErr) {
	// Jalankan Method FindByID(CustomerID) yang akan mem-passing sebuah argument berupa query string 
	// yang ada di URL Browser untuk dicocokan dengan data di database berdasarkan id(CustomerID)
	cust, err := s.repository.FindByID(CustomerID)
	// Cek Error
	if err != nil {
		return nil, err
	}

	// Setelah data-nya didapatkan, convert tipe data field status dari int ke string
	response := cust.ToDTO()

	// Jika sudah kembalikan data yang sudah di convert
	return &response, nil
}
