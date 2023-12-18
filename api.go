package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type ApiFunc func(http.ResponseWriter, *http.Request) error
type ApiServer struct {
	listenAdrr string
	store      Storage
}
type ApiError struct {
	msgError string
}

func WriteJson(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func makeHTTPHandleFunc(f ApiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			//Error handler
			WriteJson(w, http.StatusBadRequest, ApiError{msgError: err.Error()})
		}
	}
}

func newApiServer(listenAdrr string, store Storage) *ApiServer {
	return &ApiServer{
		listenAdrr: listenAdrr,
		store:      store,
	}
}

func (s *ApiServer) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/account", makeHTTPHandleFunc(s.handleAccount))
	router.HandleFunc("/account/{id}", makeHTTPHandleFunc(s.handleAccount))

	log.Println("JSON API server running on port: ", s.listenAdrr)
	http.ListenAndServe(s.listenAdrr, router)
}

func (s *ApiServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		return s.handleGetAccount(w, r)
	case "POST":
		return s.handleCreateAccount(w, r)
	case "DELETE":
		return s.handleDeleteAccount(w, r)
	}
	return fmt.Errorf("Method not allowd %s", r.Method)
}

func (s *ApiServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	//account := newAccount("Cristian", "Test")
	vars := mux.Vars(r)

	return WriteJson(w, http.StatusOK, vars)
}
func (s *ApiServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	_newAccount := new(CreateAccountRequest)

	if err := json.NewDecoder(r.Body).Decode(_newAccount); err != nil {
		return err
	}

	account := newAccount(_newAccount.FirstName, _newAccount.LastName)
	id, err := s.store.CreateAccount(account)
	if err != nil {
		return err
	}
	account.Id = id

	return WriteJson(w, http.StatusOK, account)
}
func (s *ApiServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}
