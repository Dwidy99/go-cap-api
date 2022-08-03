package app

import (
	"capi/service"
	"encoding/json"
	"encoding/xml"
	"net/http"

	"github.com/gorilla/mux"
)

type CustomerHandler struct {
	service service.CustomerService
}

func (ch *CustomerHandler) getAllCustomers(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	customers, _ := ch.service.GetAllCustomers(status)

	if r.Header.Get("Content-Type") == "application/xml" {
		w.Header().Add("Content-Type", "application/xml")
		xml.NewEncoder(w).Encode(customers)
	} else {
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(customers)
	}
}

func (ch *CustomerHandler) getCustomerByID(w http.ResponseWriter, r *http.Request) {
	// * get route variabel
	vars := mux.Vars(r)

	customerID := vars["customer_id"]

	customer, err := ch.service.GetCustomerByID(customerID)

	if err != nil {
		writeResponse(w, err.Code, err.AsMessage())
		return
	}

	// * return customer data
	writeResponse(w, http.StatusOK, customer)
}

func writeResponse(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		panic(err)
	} 
}