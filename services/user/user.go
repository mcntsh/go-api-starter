package user

import (
	"errors"
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql" // MySQL Go driver
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

// Database queries
const (
	queryInsertNewUser   = `INSERT INTO users (email, password) VALUES (:email, :password)`
	queryFindUserByEmail = `SELECT * FROM users WHERE email=?`
	queryFindUserByID    = `SELECT * FROM users WHERE id=?`
)

// Error responses
const (
	errorInvalidEmail       = "Invalid email address"
	errorEmailNotFound      = "Could not find a user with the given email address"
	errorShortPassword      = "Password must be greater than %d characters"
	errorExistingUser       = "A user with this email address has already been registered"
	errorIDNotFound         = "Could not find a user with the given ID"
	errorJWT                = "An error occured while trying to product a JSON Web Token"
	errorInvalidCredentials = "Invalid email/password combination"
)

// A User represents the database record of the user model.
type User struct {
	db *sqlx.DB

	ID       int64  `db:"id"`
	Email    string `db:"email"`
	Password string `db:"password"`
}

// A SimpleUser represents a simpler version of User without any of
// the supperfluous database fields. Used for API responses.
type SimpleUser struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
}

// A Token is a <a href="http://jwt.io">JSON Web Token</a> representation of a
// User authenticated session.
type Token struct {
	Token string `json:"token"`
}

// CreateUser creates a new User record with an email and password,
// and saves it to the database.
func (u *User) CreateUser(email, password string) error {
	u.Email = email
	u.Password = password

	// Validation
	if govalidator.IsEmail(u.Email) != true {
		return errors.New(errorInvalidEmail)
	}
	if govalidator.IsByteLength(u.Password, 8, 255) != true {
		return fmt.Errorf(errorShortPassword, 8)
	}

	// Check if there's an existing user with this email
	_, err := u.FindUserByEmail(u.Email)
	if err == nil {
		return errors.New(errorExistingUser)
	}

	// Encode the user's password
	encryptedPass, err := bcrypt.GenerateFromPassword([]byte(password), 5)
	if err != nil {
		return err
	}

	u.Password = string(encryptedPass)

	// Insert into the DB
	_, err = u.db.NamedExec(queryInsertNewUser, *u)
	if err != nil {
		return err
	}

	return err
}

// FindUserByID retrieves a User record from the database by its
// index.
func (u *User) FindUserByID(id int64) (*User, error) {
	// Find user row
	err := u.db.QueryRowx(queryFindUserByID, id).StructScan(u)
	if err != nil {
		return nil, errors.New(errorIDNotFound)
	}

	return u, err
}

// FetchUserTokenByID retrieves a user by it's primary key, and then
// encodes it into a <a href="http://jwt.io">JSON Web Token</a>.
func (u *User) FetchUserTokenByID(id int64) (*Token, error) {
	t := &Token{}

	u, err := u.FindUserByID(id)
	if err != nil {
		return nil, err
	}

	jt := jwt.New(jwt.SigningMethodHS256)
	jt.Claims["id"] = u.ID
	jt.Claims["email"] = u.Email

	t.Token, err = jt.SignedString([]byte(config.EncodingJWT))
	if err != nil {
		return nil, errors.New(errorJWT)
	}

	return t, err
}

// FindUserByEmail queries the database for a User record by its
// email/username.
func (u *User) FindUserByEmail(email string) (*User, error) {
	// Validation
	if govalidator.IsEmail(email) != true {
		return nil, errors.New(errorInvalidEmail)
	}

	// Find user row
	err := u.db.QueryRowx(queryFindUserByEmail, email).StructScan(u)
	if err != nil {
		return nil, errors.New(errorEmailNotFound)
	}

	return u, err
}

// FindUserByEmailAndPassword queries the database for a User
// record by its username/password pair.
func (u *User) FindUserByEmailAndPassword(email, password string) (*User, error) {
	u, err := u.FindUserByEmail(email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return nil, errors.New(errorInvalidCredentials)
	}

	return u, err
}

func LoadUser() (*User, error) {
	var err error

	u := &User{}

	u.db, err = sqlx.Connect(config.DBDriver, config.dbString)
	if err != nil {
		return nil, err
	}

	return u, err
}

func init() {
	config.LoadJSON()
}
