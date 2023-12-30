package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

func (s *ApiServer) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/account", makeHTTPHandleFunc(s.handleAccount))
	router.HandleFunc("/account/{id}", makeHTTPHandleFunc(s.handleAccountById))

	//User:
	router.HandleFunc("/user/register", makeHTTPHandleFunc(s.handleRegisterUser)).Methods("POST")
	router.HandleFunc("/user/login", makeHTTPHandleFunc(s.handleUserLogin)).Methods("POST")

	log.Println("JSON API server running on port", s.listenAdrr)
	http.ListenAndServe(s.listenAdrr, router)
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
			WriteJson(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

func newApiServer(listenAdrr string, store Storage) *ApiServer {
	return &ApiServer{
		listenAdrr: listenAdrr,
		store:      store,
	}
}

func (s *ApiServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		return s.handleGetAllAccount(w, r)
	case "POST":
		return s.handleCreateAccount(w, r)
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

/*--- Begin Api Methods ---*/

func (s *ApiServer) handleRegisterUser(w http.ResponseWriter, r *http.Request) error {
	_newUser := new(CreateUserRequest)

	if err := json.NewDecoder(r.Body).Decode(_newUser); err != nil {
		return err
	}

	registeredUser := new(User)
	registeredUser.CreatedAt = time.Now().UTC()
	registeredUser.Email = _newUser.Email
	registeredUser.Username = _newUser.Username

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(_newUser.Password), 8)
	if err != nil {
		return err
	}

	registeredUser.HashedPassword = string(hashedPassword)

	id, err := s.store.RegisterUser(registeredUser)
	if err != nil {
		return err
	}
	registeredUser.Id = id

	return WriteJson(w, http.StatusOK, registeredUser)
}

func (s *ApiServer) handleUserLogin(w http.ResponseWriter, r *http.Request) error {
	userVerify := new(UserLoginRequest)

	if err := json.NewDecoder(r.Body).Decode(userVerify); err != nil {
		return err
	}

	hashedPassword, err := s.store.SelectUserPassword(userVerify.Username)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(userVerify.Password)); err != nil {
		return errors.New("Contraseña incorrecta")
	}

	return WriteJson(w, http.StatusOK, "ok")
}

func (s *ApiServer) handleGetAllAccount(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.store.GetAllAccounts()

	if err != nil {
		return err
	}

	return WriteJson(w, http.StatusOK, accounts)
}

func (s *ApiServer) handleGetAccountById(w http.ResponseWriter, r *http.Request) error {
	//account := newAccount("Cristian", "Test")
	id, err := parseId(r)
	if err != nil {
		return err
	}

	account, err := s.store.GetAccountById(id)

	if err != nil {
		return err
	} else if account.Id == 0 {
		return WriteJson(w, http.StatusNotFound, "Objeto no encontrado.")
	}

	return WriteJson(w, http.StatusOK, account)
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

func (s *ApiServer) handleUpdateAccount(w http.ResponseWriter, r *http.Request) error {
	// vars := mux.Vars(r)
	account := new(Account)
	id, err := parseId(r)
	if err != nil {
		return err
	}

	if err := json.NewDecoder(r.Body).Decode(account); err != nil {
		return err
	}

	account.Id = id

	if err := s.store.UpdateAccount(account); err != nil {
		return err
	}

	return WriteJson(w, http.StatusOK, "Ok")
}

func (s *ApiServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	id, err := parseId(r)
	if err != nil {
		return err
	}

	if err := s.store.DeleteAccount(id); err != nil {
		return err
	}

	return WriteJson(w, http.StatusOK, "ok")
}

/*--- End Api Methods ---*/

func parseId(r *http.Request) (int, error) {
	id_string := mux.Vars(r)["id"]

	id, err := strconv.Atoi(id_string)

	if err != nil {
		return 0, errors.New("Formato de id inválido.")
	}
	return id, nil
}

/*--- Begin Types ---*/
type ApiFunc func(http.ResponseWriter, *http.Request) error

type ApiServer struct {
	listenAdrr string
	store      Storage
}

/*--- End Types ---*/
