package domain

import (
	"capi/errs"
	"capi/logger"
	"database/sql"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Penjelasan DB 1.
// Buat data struct yang berisi properti client dengan tipe data *sqlx.DB
type CustomerRepositoryDB struct {
	client *sqlx.DB
}

// Penjeasan DB 1.1 Buat fungsi NewCustomerRepositoryDB() untuk membuat koneksi database yang me-return
// struct CustomerRepositoryDB berisi db berupa reference, variabel db akan mengisi variabel 
// client berupa tipe data package *sqlx.BD
func NewCustomerRepositoryDB() CustomerRepositoryDB {
	connStr := "user=postgres password=d dbname=banking sslmode=disable"
	db, err := sqlx.Open("postgres", connStr) // me-return reference(&) dan error
	if err != nil {
		log.Fatal("Your Database->  ", err)
	}

	return CustomerRepositoryDB{db}
}

// Penjelasan ByID 1.
// Method FindByID(customerID string) akan mengembalikan data (belum berupa json)
// berdasarkan argument customerID berupa string
func (d CustomerRepositoryDB) FindByID(customerID string) (*Customer, *errs.AppErr) {
	// query berdasarkan id
	query := "select * from customers where customer_id = $1"

	// row := d.client.QueryRow(query, customerID)

	var c Customer
	// Get(arg1, arg2, arg3)
	// argument ke-3 akan me-replace/mem-passing $1 yang ada di argument ke-2 (berisi query string)
	// kemudian argument ke-1 &c akan menampung hasil data-nya
	err := d.client.Get(&c, query, customerID)

	// err := row.Scan(&c.ID, &c.Name, &c.DateOfBirth, &c.City, &c.ZipCode, &c.Status)
	// Cek Error Custom
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Error("Error customer Data Not Found-> "+ err.Error())
			return nil, errs.NewNotFoundError("Customer Not Found")
		} else {
			logger.Error("Error Scanning Data-> " + err.Error())
			return nil, errs.NewUnexpectedError("Unexpected Database Error")
		}
	}

	return &c, nil
}

// Penjelasan AllCustomers 1.
// Method FindAll(customerStatus string) akan mengembalikan semua data (belum berupa json)
func (d CustomerRepositoryDB) FindAll(customerStatus string) ([]Customer, *errs.AppErr) {
	var c []Customer

	// Cek apakah customerStatus tidak kosong
	if customerStatus != "" {
		// Cek apakah customerStatus "inactive"
		if customerStatus == "inactive" {
			// Set "0"
			customerStatus = "0"
		} else {
			// Set "1"
			customerStatus = "1"
		}
		// Query data
		query := "select * from customers where status = $1"

		// Select(arg1, arg2, arg3)
		// argument ke-3 akan me-replace/mem-passing $1 yang ada di argument ke-2 (berisi query string)
		// kemudian argument ke-1 &c akan menampung hasil data-nya
		err := d.client.Select(&c, query, customerStatus)
		// Cek Error
		if err != nil {
			logger.Error("Error Query Customers table" + err.Error())
			return nil, errs.NewUnexpectedError("Expected Database Error")
		}
	} else {
		query := "select * from customers"
	
		// Select(arg1, arg2, arg3)
		// argument ke-1 &c akan menampung hasil data-nya
		err := d.client.Select(&c, query)
		if err != nil {
			logger.Error("error query data to customer table "+ err.Error())
			return nil, errs.NewUnexpectedError("Unexpected DB Error")
		}
	}
	return c, nil


	// var customers []Customer
	// for rows.Next() {

	// 	var c Customer
	// 	err := rows.Scan(&c.ID, &c.Name, &c.DateOfBirth, &c.City, &c.ZipCode, &c.Status)
	// 	if err != nil {
	// 		log.Println("error scanning customer data ", err.Error())
	// 	}

	// 	customers = append(customers, c)
	// }

}