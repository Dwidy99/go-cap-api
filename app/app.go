package app

import (
	"capi/domain"
	"capi/service"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func Start() {
	// Penjelasan db 0
	// CustomerHandler (file handler.go package yg sama) adalah sebuah struct berisi service.CustomerService 
	// yang ditampung oleh variabel service untuk diisi oleh argument domain.NewCustomerRepositoryDB(). 
	// Fungsi NewCustomerRepositoryDB() berisi sebuah koneksi database yang me-return sebuah
	// struct CustomerRepositoryDB ditampung ke variabel client dengan package *sqlx.DB
	// * wiring
	ch := CustomerHandler{service.NewCustomerService(domain.NewCustomerRepositoryDB())}
	// variabel ch nantinya akan digunakan oleh method-method yang ada di argument ke-dua 
	// dari ROUTE/rute untuk keperluan query supaya tetap terkoneksi ke database

	// * create ServeRoute
	mux := mux.NewRouter()

	// Penjelasan AllCustomers 0.
	// ROUTE/rute ini untuk menangani 2 argument dari mux.HandleFunc(arg1, arg2) yang menggunakan method GET.
	// Nilai dari "/customers" yang ditangkap dari browser akan ditangani oleh ch.getAllCustomers
	mux.HandleFunc("/customers", ch.getAllCustomers).Methods("GET")

	// Penjelasan ByID 0.
	// ROUTE/rute ini untuk menangani 2 argument dari mux.HandleFunc(arg1, arg2) yang menggunakan method GET.
	// Nilai dari "/customers/{customer_id:[0-9]+}" yang ditangkap dari browser akan ditangani oleh ch.getCustomerByID
	mux.HandleFunc("/customers/{customer_id:[0-9]+}", ch.getCustomerByID).Methods("GET")

	// * starting the server
	fmt.Println("starting the server localhost:9000")
	http.ListenAndServe(":9000", mux)
}