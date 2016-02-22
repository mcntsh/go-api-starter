package main

import (
	"encoding/json"
	"github.com/gorilla/context"
	"micro-services/api"
	"net/http"
)

type registerUserBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type authUserBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func handlerGetAuthenticatedUser(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user").(*model)

	api.WriteResponse(w, &simpleUser{Id: user.Id, Email: user.Email})
}

func handlerRegisterUser(w http.ResponseWriter, r *http.Request) {
	var request *registerUserBody

	model := context.Get(r, "model").(*model)

	// Unmarshal the POST body into the body struct
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		api.WriteErrorResponse(w, r, http.StatusInternalServerError, err)
		return
	}

	// Create the new user
	err = model.CreateUser(request.Email, request.Password)
	if err != nil {
		api.WriteErrorResponse(w, r, http.StatusBadRequest, err)
		return
	}

	api.WriteResponse(w, nil)
	return
}

func handlerAuthenticateUser(w http.ResponseWriter, r *http.Request) {
	var request *authUserBody

	model := context.Get(r, "model").(*model)

	// Unmarshal the POST body into the body struct
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		api.WriteErrorResponse(w, r, http.StatusInternalServerError, err)
		return
	}

	// Authenticate the user
	user, err := model.FindUserByEmailAndPassword(request.Email, request.Password)
	if err != nil {
		api.WriteErrorResponse(w, r, http.StatusBadRequest, err)
		return
	}

	// Generate the JWT
	token, err := model.FindUserTokenById(user.Id)
	if err != nil {
		api.WriteErrorResponse(w, r, http.StatusNoContent, err)
		return
	}

	api.WriteResponse(w, token)
	return
}
