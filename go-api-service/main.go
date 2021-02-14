package main

import (
	"encoding/json"
	"fmt"
	"github.com/rs/cors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type company struct {
	CVR     int
	Name    string
	Address string
}

type companiesResponseWrapper struct {
	Companies []company
	Count     int
}

type companylist []company

var companies = companylist{
	{
		CVR:     61056416,
		Name:    "CARLSBERG A/S",
		Address: "J.C. Jacobsens Gade 1",
	},
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.Use(headerFilter)

	// Routing
	router.HandleFunc("/companies", createCompany).Methods("POST")
	router.HandleFunc("/companies", fetchAll).Methods("GET")
	router.HandleFunc("/companies/{cvr}", getCompany).Methods("GET")
	router.HandleFunc("/companies/{cvr}", updateCompany).Methods("PATCH")
	router.HandleFunc("/companies/{cvr}", deleteCompany).Methods("DELETE")

	// CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"POST", "GET", "PATCH", "DELETE"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)
	log.Fatal(http.ListenAndServe(":8080", handler))
}

func createCompany(responseWriter http.ResponseWriter, request *http.Request) {
	var newCompany company
	reqBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(responseWriter, "Bad request")
		return
	}

	json.Unmarshal(reqBody, &newCompany)

	if companyExists(newCompany.CVR) {
		responseWriter.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(responseWriter, "The company with CVR %v already exists", newCompany.CVR)
		return
	}

	companies = append(companies, newCompany)
	responseWriter.WriteHeader(http.StatusCreated)

	json.NewEncoder(responseWriter).Encode(newCompany)
}

func getCompany(responseWriter http.ResponseWriter, request *http.Request) {
	cvr := getCvrRequestParameter(responseWriter, request)

	if !companyExists(cvr) {
		responseWriter.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(responseWriter, "Not Found")
		return
	}

	for _, company := range companies {
		if company.CVR == cvr {
			json.NewEncoder(responseWriter).Encode(company)
		}
	}
}

func fetchAll(responseWriter http.ResponseWriter, request *http.Request) {
	count := len(companies)
	response := companiesResponseWrapper{companies, count}
	json.NewEncoder(responseWriter).Encode(response)
}

func updateCompany(responseWriter http.ResponseWriter, request *http.Request) {
	cvr := getCvrRequestParameter(responseWriter, request)
	var companyPatch company

	reqBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	json.Unmarshal(reqBody, &companyPatch)
	result, updated := patchCompany(cvr, companyPatch)

	if !updated {
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}
	json.NewEncoder(responseWriter).Encode(result)
}

func deleteCompany(responseWriter http.ResponseWriter, request *http.Request) {
	cvr := getCvrRequestParameter(responseWriter, request)

	for i, company := range companies {
		if company.CVR == cvr {
			companies = append(companies[:i], companies[i+1:]...)
			fmt.Fprintf(responseWriter, "The company with CVR %v has been deleted successfully", cvr)
		}
	}
}

func patchCompany(cvr int, patch company) (company, bool) {
	for i, company := range companies {
		if company.CVR == cvr {
			company.Name = patch.Name
			company.Address = patch.Address
			companies = append(companies[:i], company)
			return company, true
		}
	}
	return patch, false
}

func companyExists(cvr int) bool {
	for _, company := range companies {
		if company.CVR == cvr {
			return true
		}
	}
	return false
}

func getCvrRequestParameter(responseWriter http.ResponseWriter, request *http.Request) int {
	requestValue := mux.Vars(request)["cvr"]
	cvr, err := strconv.Atoi(requestValue)
	if err != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(responseWriter, "invalid request, cvr must be a number")
	}
	return cvr
}

func headerFilter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
