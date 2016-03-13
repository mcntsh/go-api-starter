package main

import (
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"net/http"
)

var (
	defaultChain = alice.New(
		middlewareLogging,
		middlewareJSON,
		middlewareCors,
		middlewareMountUser,
	)

	authChain = defaultChain.Extend(
		alice.New(middlewareAuth),
	)
)

func loadRouter() http.Handler {
	r := mux.NewRouter()

	// REST Handlers

	r.Methods("POST").Path("/users").Handler(defaultChain.ThenFunc(handlerRegisterUser))

	r.Methods("POST").Path("/users/auth").Handler(defaultChain.ThenFunc(handlerAuthenticateUser))
	r.Methods("GET").Path("/users/auth").Handler(authChain.ThenFunc(handlerGetAuthenticatedUser))

	// Catch-all Handler

	r.PathPrefix("/").Handler(http.DefaultServeMux)

	return r
}
