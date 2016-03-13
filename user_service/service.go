package main

import (
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/tylerb/graceful"
)

func main() {
	// Load in the config

	err := loadConfig()
	if err != nil {
		logrus.Fatal(err)
	}

	// Start the server

	httpTimeout := time.Duration(config.HTTPTimeout) * time.Second
	srv := &graceful.Server{
		Timeout: httpTimeout,
		Server: &http.Server{
			Addr:    config.HTTPAddress,
			Handler: loadRouter(),
		},
	}

	logrus.Infoln("Running HTTP server on " + config.HTTPAddress)

	err = srv.ListenAndServe()
	if err != nil {
		logrus.Fatal(err)
	}
}
