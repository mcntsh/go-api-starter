package main

import (
	"micro-services/api"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/rs/cors"
)

type userRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

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

func middlewareMountUser(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, err := loadModel()
		if err != nil {
			api.WriteErrorResponse(w, r, http.StatusInternalServerError, err)
		}

		context.Set(r, "user", u)
		h.ServeHTTP(w, r)
	})
}

func middlewareAuth(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the token
		t, err := jwt.ParseFromRequest(r, func(t *jwt.Token) (interface{}, error) {
			return []byte(config.EncodingJWT), nil
		})
		if err != nil {
			api.WriteErrorResponse(w, r, http.StatusUnauthorized, err)
			return
		}

		// Fetch the user
		u := context.Get(r, "user").(*User)

		u, err = u.FindUserByID(int64(t.Claims["id"].(float64)))
		if err != nil {
			api.WriteErrorResponse(w, r, http.StatusUnauthorized, err)
			return
		}

		context.Set(r, "user", u)

		h.ServeHTTP(w, r)
	})
}
