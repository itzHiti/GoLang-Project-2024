package model

import (
	"OCM/pkg/OCM/validator"
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type password struct {
	plaintext *string
	hash      []byte
}

type User struct {
	ID        int64    `json:"-"`
	Username  string   `json:"username"`
	Email     string   `json:"email"`
	Password  password `json:"-"`
	TokenHash string   `json:"-"`
	Activated bool     `json:"-"`
	Role      string   `json:"-"`
}

var (
	ErrDuplicateEmail    = errors.New("duplicate email")
	ErrRecordNotFound    = errors.New("Record not Found")
	ErrEditConflict      = errors.New("Error Edit Conflict")
	ErrDuplicateUsername = errors.New("Duplicated Username")
)
var AnonymousUser = &User{}

type UserModel struct {
	DB    *sql.DB
	Email string `json:"email"`
}

func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

func (u UserModel) IsAdmin(id int64) (bool, error) {
	query := `select $1 in (select user_id from admins) as isadmin`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var ans bool
	err := u.DB.QueryRowContext(ctx, query, id).Scan(&ans)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return false, ErrRecordNotFound
		default:
			return false, err
		}
	}
	return ans, nil
}

func (u UserModel) HasBan(id int64) (bool, error) {
	query := `
	select bans.id
	from bans
	inner join users on users.id=bans.user_id
	where bans.user_id = $1
	and bans.expiry > $2`

	var ban_id int64
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := u.DB.QueryRowContext(ctx, query, id, time.Now()).Scan(
		&ban_id,
	)

	query = `
	delete from bans
       where expiry<$1`
	_, _ = u.DB.Exec(query, time.Now())
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

func (u UserModel) GetRole(id int64) (string, error) {
	isadmin, _ := u.IsAdmin(id)
	if isadmin {
		return "admin", nil
	}

	ban, _ := u.HasBan(id)
	if ban {
		return "banned_user", nil
	}
	return "user", nil
}

func (u UserModel) Insert(user *User) error {
	query := `
	INSERT INTO users (username,email, password, token_hash)
	VALUES ($1, $2, $3, $4)
	RETURNING id`
	args := []interface{}{user.Username, user.Email, user.Password.hash, user.TokenHash}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := u.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID)

	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_username_key"`:
			return ErrDuplicateUsername
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}
	return nil
}

func (m UserModel) GetByUsername(username string) (*User, error) {
	query := `
	select users.id, username, email,password, activated, count(rating) as reviews
	from users left join reviews
	on reviews.user_id=users.id
	WHERE username = lower($1)
	group by users.id;`
	var user User
	err := m.DB.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password.hash,
		&user.Activated,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}

func (m UserModel) GetByEmail(email string) (*User, error) {
	query := `
	SELECT id, username, email, password, activated
	FROM users
	WHERE email = $1`
	var user User
	err := m.DB.QueryRow(query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password.hash,
		&user.Activated,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}

func (u UserModel) GetByVerificationCode(plaintext string) (*User, error) {
	hash := sha256.Sum256([]byte(plaintext))
	log.Printf("Debug: Hash generated for plaintext '%s': %x", plaintext, hash)

	query := `
    SELECT users.id, users.username, users.email, users.password, users.activated, users.token_hash
    FROM users
    INNER JOIN verifications
    ON users.id = verifications.user_id
    WHERE verifications.code = $1`

	args := []interface{}{hash[:]}
	log.Printf("Debug: Query args set: %x", plaintext)

	var user User
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := u.DB.QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password.hash,
		&user.Activated,
		&user.TokenHash,
	)

	if err != nil {
		log.Printf("Error: QueryRowContext failed with %v", err)
		switch {
		case errors.Is(err, sql.ErrNoRows):
			log.Printf("Debug: No rows found for given hash")
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	log.Printf("Debug: User retrieved with ID %d", user.ID)
	return &user, nil
}

func (m UserModel) Update(user *User) error {
	query := `
	UPDATE users
	SET username = $1, email = $2, password = $3, activated = $4, token_hash=$5
	WHERE id = $6 
	RETURNING id`
	args := []interface{}{
		user.Username,
		user.Email,
		user.Password.hash,
		user.Activated,
		user.TokenHash,
		user.ID,
	}
	err := m.DB.QueryRow(query, args...).Scan()
	if errors.Is(err, sql.ErrNoRows) {
		return ErrEditConflict
	}
	return nil
}

func (u UserModel) Delete(username string) error {
	query := `
		DELETE FROM users
		WHERE username = $1`
	result, err := u.DB.Exec(query, username)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}
	p.plaintext = &plaintextPassword
	p.hash = hash
	return nil
}
func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}
func ValidateUsername(v *validator.Validator, username string) {
	v.Check(username != "", "username", "must be provided")
	v.Check(validator.Matches(username, validator.UsernameRX), "username", "must be a valid username")
}
func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}

func ValidateEmailOrUsername(v *validator.Validator, username string, email string) {
	v1 := validator.New()
	ValidateUsername(v1, username)
	v2 := validator.New()
	ValidateEmail(v2, email)
	v.Check(v1.Valid() || v2.Valid(), "error", "valid username or email must be provided")
}

func ValidateUser(v *validator.Validator, user *User) {
	ValidateUsername(v, user.Username)
	ValidateEmail(v, user.Email)
	if user.Password.plaintext != nil {
		ValidatePasswordPlaintext(v, *user.Password.plaintext)
	}
	if user.Password.hash == nil {
		panic("missing password hash for user")
	}
}

func (u UserModel) ActivateUser(userID int64) error {
	query := `UPDATE users SET activated = true WHERE id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := u.DB.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("No rows affected")
	}

	return nil
}
