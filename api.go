package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func makeHTTPHanldeFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

func NewAPIServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/account", makeHTTPHanldeFunc(s.handleAccount))
	router.HandleFunc("/account/{id}", makeHTTPHanldeFunc(s.handleGetAccountById))
	router.HandleFunc("/account/{id}", makeHTTPHanldeFunc(s.handleUpdateAccount)).Methods("PUT")
	router.HandleFunc("/accounts", makeHTTPHanldeFunc(s.handleGetAccount))
	router.HandleFunc("/transfer", makeHTTPHanldeFunc(s.handleTransfer)).Methods("POST")

	log.Println("JSON API Server Running on port ", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetAccount(w, r)
	}
	if r.Method == "POST" {
		return s.handleCreateAccount(w, r)
	}
	return fmt.Errorf("method not allowed %s", r.Method)
}

// GET /account
func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.store.GetAccounts()

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, accounts)
}

func (s *APIServer) handleGetAccountById(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		id, err := getID(r)

		if err != nil {
			return fmt.Errorf("invalid id given")
		}

		account, err := s.store.GetAccountByID(id)

		if err != nil {
			return err
		}

		return WriteJSON(w, http.StatusOK, account)
	case "DELETE":
		return s.handleDeleteAccount(w, r)
	}

	return fmt.Errorf("method not allowed")
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	createAccountReq := new(CreateAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(createAccountReq); err != nil {
		return err
	}

	account := NewAccount(createAccountReq.FirstName, createAccountReq.LastName)
	if err := s.store.CreateAccount(account); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusCreated, account)
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}

	if err := s.store.DeleteAccount(id); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, map[string]int{"deleted": id})
}

func (s *APIServer) handleUpdateAccount(w http.ResponseWriter, r *http.Request) error {
	updateAccountReq := new(UpdateAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(updateAccountReq); err != nil {
		return err
	}

	id, err := strconv.Atoi(mux.Vars(r)["id"])

	if err != nil {
		return err
	}

	account, err := s.store.GetAccountByID(id)

	if err != nil {
		return err
	}

	account.FirstName = updateAccountReq.FirstName
	account.LastName = updateAccountReq.LastName

	if err := s.store.UpdateAccount(account); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, "updated")
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	transferReq := new(TransferRequest)
	if err := json.NewDecoder(r.Body).Decode(transferReq); err != nil {
		return err
	}
	defer r.Body.Close()

	return WriteJSON(w, http.StatusOK, transferReq)
}

func getID(r *http.Request) (int, error) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)

	if err != nil {
		return id, fmt.Errorf("invalid id given %s", idStr)
	}

	return id, nil
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type APIServer struct {
	listenAddr string
	store      Storage
}

type ApiError struct {
	Error string `json:"error"`
}
