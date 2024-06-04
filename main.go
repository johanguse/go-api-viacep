package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	serverPortNumber = "8080"
	apiURLViaCep     = "https://viacep.com.br/ws/%s/json"
)

func viaCepAPIHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	cep := vars["cep"]

	if len(cep) != 8 || cep == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "CEP must have 8 digits")
		return
	}

	for _, char := range cep {
		if char < '0' || char > '9' {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "CEP can be only numbers")
			return
		}
	}

	response, err := http.Get(fmt.Sprintf(apiURLViaCep, cep))
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(response.Body).Decode(&result)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "OK")
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/healthcheck", healthcheckHandler).Methods("GET")
	router.HandleFunc("/cep/{cep}", viaCepAPIHandler).Methods("GET")

	fmt.Println("Server is starting at port", serverPortNumber)

	log.Fatal(http.ListenAndServe(":"+serverPortNumber, router))
}
