package app

import (
	"capi/service"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type CustomerHandler struct {
	service service.CustomerService
}

// Penjelasan AllCustomers 3.
// Method getAllCustomers() akan menampilkan respon berupa data json
func (ch *CustomerHandler) getAllCustomers(w http.ResponseWriter, r *http.Request) {
	// Ambil request yang ada di URL browser berdasarkan(?) "status"-nya 
	// (localhost:9000/customers?status=active)
	customerStatus := r.URL.Query().Get("status")

	// Berupa semua data customers ber-tipe array / slice
	customers, err := ch.service.GetAllCustomers(customerStatus)

	// Cek error
	if err !=nil {
		// Jika ada error
		writeResponse(w, err.Code, err.AsMessage())
		return
	} else {
		// Jika tidak ada error tampilkan semua datanya berupa json
		fmt.Fprintln(w, "Customers Successfully Display!")
		
		// writeResponse(w, http.StatusOK, customers) 
		// variabel customers (argument ke-3) yang berupa data bukan json akan di convert 
		// menjadi data ber-tipe json
		writeResponse(w, http.StatusOK, customers)
	}

	// if r.Header.Get("Content-Type") == "application/xml" {
	// 	w.Header().Add("Content-Type", "application/xml")
	// 	xml.NewEncoder(w).Encode(customers)
	// } else {
	// 	w.Header().Add("Content-Type", "application/json")
	// 	json.NewEncoder(w).Encode(customers)
	// }
}

// Penjelasan ByID 3.
// Method getCustomerByID() akan menampilkan respon berupa data json
func (ch *CustomerHandler) getCustomerByID(w http.ResponseWriter, r *http.Request) {
	// * get route variabel
	// mengambil nilai parameter yang ada di URL browser
	vars := mux.Vars(r) // mengembalikan map[customer_id: "apapun"]

	customerID := vars["customer_id"] // keluarkan value id dari map

	// ch.service.GetCustomerByID(customerID) berupa data yang sudah di-convert (hanya field status)
	// (ada di package service file customerService.go)
	customer, err := ch.service.GetCustomerByID(customerID)

	// Cek Error
	if err != nil {
		writeResponse(w, err.Code, err.AsMessage())
		return
	}

	// writeResponse(w, http.StatusOK, customer) 
	// variabel customer (argument ke-3) yang berupa data bukan json akan di convert menjadi data ber-tipe json
	writeResponse(w, http.StatusOK, customer)
}

func writeResponse(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		panic(err)
	} 
}