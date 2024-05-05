package model

import (
	"database/sql"
	"errors"
	"time"
)

type BanModel struct {
	DB *sql.DB
}
type Ban struct {
	Id     int64     `json:"id"`
	UserId int64     `json:"user_id"`
	Expiry time.Time `json:"expiry"`
}

var (
	ErrDuplicateBan = errors.New("duplicate ban")
)

func (b BanModel) Insert(id int64, days int) (*Ban, error) {
	query := `
	insert into bans (user_id, expiry)
	values ($1, $2)
	returning id, user_id, expiry`

	args := []interface{}{id, time.Now().Add(time.Duration(24 * time.Hour))}
	var ban Ban

	err := b.DB.QueryRow(query, args...).Scan(
		&ban.Id,
		&ban.UserId,
		&ban.Expiry,
	)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "unique_ban"`:
			return nil, ErrDuplicateBan
		default:
			return nil, err
		}
	}
	return &ban, nil
}
