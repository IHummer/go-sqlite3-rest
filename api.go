package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

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
	router.HandleFunc("/account/{id}", makeHTTPHandleFunc(s.handleAccountById))

	log.Println("JSON API server running on port", s.listenAdrr)
	http.ListenAndServe(s.listenAdrr, router)
}

func (s *ApiServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		return s.handleGetAllAccount(w, r)
	case "POST":
		return s.handleCreateAccount(w, r)
		// case "DELETE":
		// 	return s.handleDeleteAccount(w, r)
	}
	return fmt.Errorf("Method not allowd %s", r.Method)
}

func (s *ApiServer) handleAccountById(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		return s.handleGetAccountById(w, r)
	case "POST":
		return s.handleUpdateAccount(w, r)
	case "DELETE":
		return s.handleDeleteAccount(w, r)
	}
	return fmt.Errorf("Method not allowed %s", r.Method)
}

func (s *ApiServer) handleGetAllAccount(w http.ResponseWriter, r *http.Request) error {
	//account := newAccount("Cristian", "Test")
	// vars := mux.Vars(r)

	accounts, err := s.store.GetAllAccounts()

	if err != nil {
		return WriteJson(w, http.StatusInternalServerError, "Error"+err.Error())
	}

	return WriteJson(w, http.StatusOK, accounts)
}

func (s *ApiServer) handleGetAccountById(w http.ResponseWriter, r *http.Request) error {
	//account := newAccount("Cristian", "Test")
	id_string := mux.Vars(r)["id"]

	id, err := strconv.Atoi(id_string)

	if err != nil {
		return WriteJson(w, http.StatusBadRequest, "Formato de id inválido.")
	}

	account, err := s.store.GetAccountById(id)

	if err != nil {
		return WriteJson(w, http.StatusInternalServerError, "Error: "+err.Error())
	} else if account.Id == 0 {
		return WriteJson(w, http.StatusNotFound, "Objeto no encontrado.")
	}

	return WriteJson(w, http.StatusOK, account)
}

func (s *ApiServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	_newAccount := new(CreateAccountRequest)

	if err := json.NewDecoder(r.Body).Decode(_newAccount); err != nil {
		return WriteJson(w, http.StatusBadRequest, err.Error())
	}

	account := newAccount(_newAccount.FirstName, _newAccount.LastName)
	id, err := s.store.CreateAccount(account)
	if err != nil {
		return WriteJson(w, http.StatusInternalServerError, err.Error())
	}
	account.Id = id

	return WriteJson(w, http.StatusOK, account)
}

func (s *ApiServer) handleUpdateAccount(w http.ResponseWriter, r *http.Request) error {
	// vars := mux.Vars(r)
	account := new(Account)
	id_string := mux.Vars(r)["id"]

	id, err := strconv.Atoi(id_string)

	if err != nil {
		return WriteJson(w, http.StatusBadRequest, "Formato de id inválido.")
	}

	if err := json.NewDecoder(r.Body).Decode(account); err != nil {
		return WriteJson(w, http.StatusBadRequest, err.Error())
	}

	account.Id = id

	err = s.store.UpdateAccount(account)

	if err != nil {
		return WriteJson(w, http.StatusInternalServerError, err.Error())
	}

	return WriteJson(w, http.StatusOK, "Ok")
}

func (s *ApiServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	id_string := mux.Vars(r)["id"]

	id, err := strconv.Atoi(id_string)

	if err != nil {
		return WriteJson(w, http.StatusBadRequest, "Formato de id inválido.")
	}

	err = s.store.DeleteAccount(id)

	if err != nil {
		return WriteJson(w, http.StatusBadRequest, "Error:"+err.Error())
	}

	return WriteJson(w, http.StatusOK, "ok")
}
