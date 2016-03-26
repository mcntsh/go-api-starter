package main

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/rs/cors"
)

func middlewareCors(h http.Handler) http.Handler {
	c := cors.New(cors.Options{
		AllowedOrigins:     []string{"*"},
		AllowedMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowedHeaders:     []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "X-Auth-Token"},
		OptionsPassthrough: true,
		AllowCredentials:   true,
	})

	return c.Handler(h)
}

func middlewareLogging(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logrus.WithFields(logrus.Fields{
			"method": r.Method,
			"url":    r.URL,
		}).Info("HTTP Request to handler")

		h.ServeHTTP(w, r)
	})
}

func middlewareJSON(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		h.ServeHTTP(w, r)
	})
}
