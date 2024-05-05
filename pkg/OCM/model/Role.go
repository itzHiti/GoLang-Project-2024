package model

import "database/sql"

type RoleModel struct {
	DB *sql.DB
}

type Roles []string

func (p Roles) Include(role string) bool {
	if p == nil {
		return true
	}
	for i := range p {
		if role == p[i] {
			return true
		}
	}
	return false
}

func (db *RoleModel) GetUserRole(userID int) (string, error) {
	var role string
	err := db.DB.QueryRow("SELECT role FROM users WHERE id = ?", userID).Scan(&role)
	if err != nil {
		return "", err
	}
	return role, nil
}
