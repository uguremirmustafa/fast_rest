package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type APIServer struct {
	listenAddr string
	store      Storage
}

func NewAPIServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *APIServer) Run() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /account/{id}", makeHTTPHandleFunc(s.handleGetAccount))
	mux.HandleFunc("GET /account", makeHTTPHandleFunc(s.handleGetAccounts))
	mux.HandleFunc("POST /account", makeHTTPHandleFunc(s.handleCreateAccount))

	log.Printf("JSON API Server running on port: %s \n", s.listenAddr)

	http.ListenAndServe(s.listenAddr, mux)
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	// account := NewAccount("Ugur", "Emirmustafa")
	id := r.PathValue("id")
	return WriteJSON(w, http.StatusOK, id)
}

func (s *APIServer) handleGetAccounts(w http.ResponseWriter, r *http.Request) error {

	accounts, err := s.store.GetAccounts()
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, accounts)
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	accountReq := new(CreateAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(accountReq); err != nil {
		return err
	}

	account := NewAccount(accountReq.FirstName, accountReq.LastName)

	if err := s.store.CreateAccount(account); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, account)
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			// handle error here so that
			// we don't have to handle error in the handler itself
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}

	}
}
