package user

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

func helperFetchAuthed(r *http.Request) (*User, error) {
	// Validate the token
	t, err := jwt.ParseFromRequest(r, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.EncodingJWT), nil
	})
	if err != nil {
		return nil, err
	}

	// Fetch the user
	u, err := LoadModel()
	if err != nil {
		return nil, err
	}

	u, err = u.FindUserByID(int64(t.Claims["id"].(float64)))
	if err != nil {
		return nil, err
	}

	return u, err
}
