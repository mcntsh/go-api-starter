package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/mcntsh/go-api-starter/services/post"
	"github.com/mcntsh/go-api-starter/services/user"
)

var (
	apiChain = alice.New(
		middlewareLogging,
		middlewareJSON,
		middlewareCors,
	)
)

func routerUser(r *mux.Router) {
	ur := r.PathPrefix("/users").Subrouter()

	ur.Methods("POST").Path("/").Handler(apiChain.ThenFunc(user.ActionRegisterUser))

	ur.Methods("POST").Path("/auth").Handler(apiChain.ThenFunc(user.ActionAuthenticateUser))
	ur.Methods("GET").Path("/auth").Handler(apiChain.ThenFunc(user.ActionGetAuthenticatedUser))
}

func routerPost(r *mux.Router) {
	ur := r.PathPrefix("/posts").Subrouter()

	ur.Methods("GET").Path("/").Handler(apiChain.ThenFunc(post.ActionTest))
}

// LoadRouter creates a new gorilla sub-router and defines
// the service endpoints and their handlers.
func LoadRouters() http.Handler {
	r := mux.NewRouter().StrictSlash(true)

	routerUser(r)
	routerPost(r)

	r.PathPrefix("/").Handler(http.DefaultServeMux)

	return r
}
