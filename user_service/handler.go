package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/context"
	"github.com/mcntsh/go-api"
)

/* Request body schema
What router handlers expect for certain endpoints. */
type registerUserBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type authUserBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func handlerGetAuthenticatedUser(w http.ResponseWriter, r *http.Request) {
	u := context.Get(r, "user").(*User)

	api.WriteResponse(w, &SimpleUser{ID: u.ID, Email: u.Email})
}

func handlerRegisterUser(w http.ResponseWriter, r *http.Request) {
	var body *registerUserBody

	u := context.Get(r, "user").(*User)

	// Unmarshal the POST body into the body struct
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		api.WriteErrorResponse(w, r, http.StatusInternalServerError, err)
		return
	}

	// Create the new user
	err = u.CreateUser(body.Email, body.Password)
	if err != nil {
		api.WriteErrorResponse(w, r, http.StatusBadRequest, err)
		return
	}

	api.WriteResponse(w, nil)
	return
}

func handlerAuthenticateUser(w http.ResponseWriter, r *http.Request) {
	var body *authUserBody

	u := context.Get(r, "user").(*User)

	// Unmarshal the POST body into the body struct
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		api.WriteErrorResponse(w, r, http.StatusInternalServerError, err)
		return
	}

	// Authenticate the user
	u, err = u.FindUserByEmailAndPassword(body.Email, body.Password)
	if err != nil {
		api.WriteErrorResponse(w, r, http.StatusBadRequest, err)
		return
	}

	// Generate the JWT
	t, err := u.FetchUserTokenByID(u.ID)
	if err != nil {
		api.WriteErrorResponse(w, r, http.StatusNoContent, err)
		return
	}

	api.WriteResponse(w, t)
	return
}
