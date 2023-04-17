package dba

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
)

type UserRole string

const (
	UserRoleADMIN UserRole = "ADMIN"
	UserRoleUSER  UserRole = "USER"
)

func (e *UserRole) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = UserRole(s)
	case string:
		*e = UserRole(s)
	default:
		return fmt.Errorf("unsupported scan type for UserRole: %T", src)
	}
	return nil
}

type NullUserRole struct {
	UserRole UserRole
	Valid    bool // Valid is true if UserRole is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullUserRole) Scan(value interface{}) error {
	if value == nil {
		ns.UserRole, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.UserRole.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullUserRole) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.UserRole), nil
}

type Message struct {
	ID        string         `json:"id"`
	CreatedAt int64          `json:"createdAt"`
	UpdatedAt int64          `json:"updatedAt"`
	Content   sql.NullString `json:"content"`
	ImageUrl  sql.NullString `json:"imageUrl"`
	UserId    string         `json:"userId"`
}

type User struct {
	ID        string   `json:"id"`
	CreatedAt int64    `json:"createdAt"`
	UpdatedAt int64    `json:"updatedAt"`
	Role      UserRole `json:"role"`
	Username  string   `json:"username"`
	Password  string   `json:"password"`
}
