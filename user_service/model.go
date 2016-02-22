package main

import (
	"errors"
	"github.com/asaskevich/govalidator"
	"github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

const (
	INSERT_NEW_USER    = `INSERT INTO users (email, password) VALUES (:email, :password)`
	FIND_USER_BY_EMAIL = `SELECT * FROM users WHERE email=?`
	FIND_USER_BY_ID    = `SELECT * FROM users WHERE id=?`
)

type model struct {
	db *sqlx.DB `json:"-"`

	Id       int64  `db:"id"`
	Email    string `db:"email"`
	Password string `db:"password"`
}

type simpleUser struct {
	Id    int64  `json:"id"`
	Email string `json:"email"`
}

type token struct {
	Token string `json:"token"`
}

func (m *model) CreateUser(email, password string) error {
	user := &model{}

	user.Email = email
	user.Password = password

	// Validation
	if govalidator.IsEmail(user.Email) != true {
		return errors.New("Invalid email address")
	}
	if govalidator.IsByteLength(user.Password, 8, 255) != true {
		return errors.New("Password must be greater than 8 characters")
	}

	// Check if there's an existing user with this email
	_, err := m.FindUserByEmail(user.Email)
	if err == nil {
		return errors.New("A user with this email address has already been registered")
	}

	// Encode the user's password
	encryptedPass, err := bcrypt.GenerateFromPassword([]byte(password), 5)
	if err != nil {
		return err
	}

	user.Password = string(encryptedPass)

	// Insert into the DB
	_, err = m.db.NamedExec(INSERT_NEW_USER, *user)
	if err != nil {
		return err
	}

	return err
}

func (m *model) FindUserById(id int64) (*model, error) {
	user := &model{}

	// Find user row
	err := m.db.QueryRowx(FIND_USER_BY_ID, id).StructScan(user)
	if err != nil {
		return nil, errors.New("Could not find a user with the given ID")
	}

	return user, err
}

func (m *model) FindUserTokenById(id int64) (*token, error) {
	tok := &token{}

	user, err := m.FindUserById(id)
	if err != nil {
		return nil, err
	}

	t := jwt.New(jwt.SigningMethodHS256)
	t.Claims["id"] = user.Id
	t.Claims["email"] = user.Email

	tok.Token, err = t.SignedString([]byte(config.EncodingJWT))
	if err != nil {
		return nil, errors.New("An error occured while trying to product a JSON Web Token")
	}

	return tok, err
}

func (m *model) FindUserByEmail(email string) (*model, error) {
	user := &model{}

	// Validation
	if govalidator.IsEmail(email) != true {
		return nil, errors.New("Invalid email address")
	}

	// Find user row
	err := m.db.QueryRowx(FIND_USER_BY_EMAIL, email).StructScan(user)
	if err != nil {
		return nil, errors.New("Could not find a user with the given email address")
	}

	return user, err
}

func (m *model) FindUserByEmailAndPassword(email, password string) (*model, error) {
	user, err := m.FindUserByEmail(email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("Invalid email/password combination")
	}

	return user, err
}

func loadModel() (*model, error) {
	var err error

	user := &model{}

	user.db, err = sqlx.Connect(config.DBDriver, config.DBString)
	if err != nil {
		return nil, err
	}

	return user, err
}
