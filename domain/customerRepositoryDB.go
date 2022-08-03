package domain

import (
	"capi/errs"
	"capi/logger"
	"database/sql"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type CustomerRepositoryDB struct {
	db *sqlx.DB
}

func NewCustomerRepositoryDB(client *sqlx.DB) CustomerRepositoryDB {
	connStr := "user=postgres password=d dbname=banking sslmode=disable"
	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Your Database->  ", err)
	}

	return CustomerRepositoryDB{db}
}

func (d CustomerRepositoryDB) FindByID(customerID string) (*Customer, *errs.AppErr) {
	query := "select * from customers where customer_id = $1"

	// row := d.client.QueryRow(query, customerID)

	var c Customer
	err := d.db.Get(&c, query, customerID)
	// err := row.Scan(&c.ID, &c.Name, &c.DateOfBirth, &c.City, &c.ZipCode, &c.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Error("error customer data not found "+ err.Error())
			return nil, errs.NewNotFoundError("Customer Not Found")
		} else {
			logger.Error("error scanning data " + err.Error())
			return nil, errs.NewUnexpectedError("Unexpected Database Error")
		}
	}

	return &c, nil
}

func (d CustomerRepositoryDB) FindAll() ([]Customer, error) {
	query := "select * from customers"

	rows, err := d.db.Query(query)
	if err != nil {
		log.Println("error query data to customer table ", err.Error())
		return nil, err
	}

	var customers []Customer
	for rows.Next() {

		var c Customer
		err := rows.Scan(&c.ID, &c.Name, &c.DateOfBirth, &c.City, &c.ZipCode, &c.Status)
		if err != nil {
			log.Println("error scanning customer data ", err.Error())
		}

		customers = append(customers, c)
	}

	return customers, nil
}