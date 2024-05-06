package model

import (
	"OCM/pkg/OCM/validator"
	"context"
	"crypto/rand"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func (m UserModel) Register(user *UserModel) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*user.Password.plaintext), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	activationToken := generateActivationToken()
	query := `
	INSERT INTO users (name, email, password_hash, activation_token)
	VALUES ($1, $2, $3, $4)
	RETURNING id`
	args := []interface{}{user.Name, user.Email, string(hashedPassword), activationToken}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err = m.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID)
	if err != nil {
		return err
	}
	return nil
}

func generateActivationToken() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatalf("Error generating token: %v", err)
	}
	return fmt.Sprintf("%x", b)
}

func (m UserModel) ActivateUser(token string) error {
	query := `
UPDATE users SET activated = TRUE, activation_token = NULL
WHERE activation_token = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	res, err := m.DB.ExecContext(ctx, query, token)
	if err != nil {
		return err
	}
	if count, err := res.RowsAffected(); err != nil || count != 1 {
		return errors.New("invalid or expired activation token")
	}
	return nil
}

func (u *UserModel) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}

	hash := string(bytes)
	u.Password = Password{hash: &hash}
	return nil
}

func (u *UserModel) CheckPassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(*u.Password.hash), []byte(password))
	return err
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}

func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}
func ValidateUser(v *validator.Validator, user *UserModel) {
	v.Check(user.Name != "", "name", "must be provided")
	v.Check(len(user.Name) <= 500, "name", "must not be more than 500 bytes long")

	ValidateEmail(v, user.Email)

	if user.Password.plaintext != nil {
		ValidatePasswordPlaintext(v, *user.Password.plaintext)
	}

	if user.Password.hash == nil {
		panic("missing password hash for user")
	}
}

// inserting, updating, getbyemail()
func (m UserModel) Insert(user *UserModel) error {
	// Assume a new user is not activated by default
	user.Activated = false
	query := `
INSERT INTO users (name, email, password_hash, activated)
VALUES ($1, $2, $3, $4)
RETURNING id, created_at, version`
	args := []interface{}{user.Name, user.Email, *user.Password.hash, user.Activated}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt, &user.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}
	return nil
}

func (m UserModel) GetByEmail(email string) (*UserModel, error) {
	query := `
SELECT id, created_at, name, email, password_hash, activated, version
FROM users
WHERE email = $1`
	var user UserModel
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Name,
		&user.Email,
		&user.Password.hash,
		&user.Activated,
		&user.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, err
		default:
			return nil, err
		}
	}
	return &user, nil
}

func (m UserModel) Update(user *UserModel) error {
	query := `
UPDATE users
SET name = $1, email = $2, password_hash = $3, activated = $4, version = version + 1
WHERE id = $5 AND version = $6
RETURNING version`
	args := []interface{}{
		user.Name,
		user.Email,
		user.Password.hash,
		user.Activated,
		user.ID,
		user.Version,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		case errors.Is(err, sql.ErrNoRows):
			return err
		default:
			return err
		}
	}
	return nil
}
