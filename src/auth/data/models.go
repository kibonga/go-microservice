package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

const dbTimeout = time.Second * 3

var db *sql.DB

// Models is a structure that wraps all the available types for the app
type Models struct {
	User User
}

// User is a structure that represents user from the database
type User struct {
	Id        int    `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	// "-", this means that this field should not be encoded, security reasons
	Password  string    `json:"-"`
	Active    int       `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// New creates a new instance of the data package
// Returns Models type, which contains all the available types for the app
func New(dbPool *sql.DB) Models {
	db = dbPool

	return Models{
		User: User{},
	}
}

// GetAll gets all the users from the database
// Returns a slice of users, sorted by the last name
func (u *User) GetAll() ([]*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, email, first_name, last_name, password, user_active, created_at, updated_at 
			from users 
			order by last_name`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	var users []*User
	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.Id,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.Password,
			&user.Active,
			&user.CreatedAt,
			&user.UpdatedAt)
		if err != nil {
			log.Println("Error scanning", err)
			return nil, err
		}

		users = append(users, &user)
	}

	return users, nil
}

// Returns a single user by Email
func (u *User) GetByEmail(email string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, email, first_name, last_name, password, user_active, created_at, updated_at
			from users where email = $1`

	var user User
	row := db.QueryRowContext(ctx, query, email)

	err := row.Scan(
		&user.Id,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt)

	if err != nil {
		fmt.Println("Error scanning, there is no user with given email")
		return nil, err
	}

	return &user, nil
}

// Returns a single user by Id
func (u *User) GetById(id int) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, email, first_name, last_name, password, user_active, created_at, updated_at
			from users where id = $1`

	var user User
	row := db.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&user.Id,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt)

	if err != nil {
		fmt.Println("Error scanning, there is no user with given id")
		return nil, err
	}

	return &user, nil
}

// Updates user
func (u *User) Update() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	cancel()

	query := `update users set
				email = $1,
				first_name = $2,
				last_name = $3,
				user_active = $4,
				updated_at = $5,
				where id = $6`

	_, err := db.ExecContext(ctx, query, u.Email, u.FirstName, u.LastName, u.Active, u.UpdatedAt)
	if err != nil {
		fmt.Println("Error executing update")
		return err
	}

	return nil
}

// Deletes a single user, by User.Id
func (u *User) Delete() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `delete from users where id = $1`

	_, err := db.ExecContext(ctx, query, u.Id)
	if err != nil {
		fmt.Println("Error executing delete")
		return err
	}

	return nil
}

// Deletes a single user by id
func (u *User) DeletById(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `delete from users where id = $1`
	_, err := db.ExecContext(ctx, query, id)
	if err != nil {
		fmt.Println("Error executing delete by id")
		return err
	}

	return nil
}

// Inserts a new user
// Returns Id of the newly created user
func (u *User) Insert() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), 12)
	if err != nil {
		return 0, err
	}

	var id int
	query := `insert into users (email, first_name, last_name, password, user_active, created_at, updated_at)
			values($1, $2, $3, $4, $5, $6, $7) returning id`

	err = db.QueryRowContext(ctx, query, u.Email, u.FirstName, u.LastName, hashedPassword, u.Active, time.Now(), time.Now()).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

// Resets users password
func (u *User) ResetPassword(password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), 12)
	if err != nil {
		return err
	}

	query := `update users set password = $1 where id = $2`

	_, err = db.ExecContext(ctx, query, hashedPassword, u.Id)
	if err != nil {
		return err
	}

	return nil
}

// Check if users password matches bcrypt encrypted password
func (u *User) PasswordMatches(plainTextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainTextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			// invalid password
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}
