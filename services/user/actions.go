package user

import (
	"encoding/json"
	"net/http"

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

func ActionGetAuthenticatedUser(w http.ResponseWriter, r *http.Request) {
	u, err := helperFetchAuthed(r)
	if err != nil {
		api.WriteErrorResponse(w, r, http.StatusUnauthorized, err)
		return
	}

	api.WriteResponse(w, &SimpleUser{ID: u.ID, Email: u.Email})
}

func ActionRegisterUser(w http.ResponseWriter, r *http.Request) {
	var body *registerUserBody

	u, err := LoadUser()
	if err != nil {
		api.WriteErrorResponse(w, r, http.StatusInternalServerError, err)
		return
	}

	// Unmarshal the POST body into the body struct
	err = json.NewDecoder(r.Body).Decode(&body)
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

func ActionAuthenticateUser(w http.ResponseWriter, r *http.Request) {
	var body *authUserBody

	u, err := LoadUser()
	if err != nil {
		api.WriteErrorResponse(w, r, http.StatusUnauthorized, err)
		return
	}

	// Unmarshal the POST body into the body struct
	err = json.NewDecoder(r.Body).Decode(&body)
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
