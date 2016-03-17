package userserv

import (
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
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

// LoadRouter creates a new gorilla sub-router and defines
// the service endpoints and their handlers.
func (s *Service) LoadRouter(path string, r *mux.Router) {
	sr := r.PathPrefix(path).Subrouter()

	// REST Handlers
	sr.Methods("POST").Path("/").Handler(defaultChain.ThenFunc(handlerRegisterUser))

	sr.Methods("POST").Path("/auth").Handler(defaultChain.ThenFunc(handlerAuthenticateUser))
	sr.Methods("GET").Path("/auth").Handler(authChain.ThenFunc(handlerGetAuthenticatedUser))

	return
}
