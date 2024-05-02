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
