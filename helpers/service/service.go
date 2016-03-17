package service

import (
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/tylerb/graceful"
)

// Servicer describes components that can be used as application
// services.
type Servicer interface {
	LoadRouter(string, *mux.Router)
}

// App represents an application server which has a router
// and sub-routers for each service.
type App struct {
	router *mux.Router
}

// NewService appends the service to the application server
// by mounting the Mux router it exports.
func (a *App) NewService(path string, serv Servicer) {
	serv.LoadRouter(path, a.router)
}

// Listen mounts all of the service routers to the main
// router and creates a graceful HTTP server listening on
// a specific port.
func (a *App) Listen(address string) {
	var err error

	// Start the server
	httpTimeout := time.Duration(5) * time.Second
	srv := &graceful.Server{
		Timeout: httpTimeout,
		Server: &http.Server{
			Addr:    address,
			Handler: a.router,
		},
	}

	logrus.Infoln("Running HTTP server on " + address)

	err = srv.ListenAndServe()
	if err != nil {
		logrus.Fatal(err)
	}
}

// NewApp creates a new application server instance with an
// attached router.
func NewApp() *App {
	a := &App{
		router: mux.NewRouter(),
	}

	return a
}
