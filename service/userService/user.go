package userService

import (
	"database/sql"
	"strings"
	"time"
)

type User struct {
	UserID   string
	Username string
	// Password keeps the plain password when creating or updating a user.
	// Important: It will never be returned to the clients.
	Password   string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	LastActive sql.NullTime
	Email      sql.NullString
	Phone      sql.NullString
	Status     int
}

func (u User) Validate() error {
	if u.Username == "" {
		return NewMissingFieldError("username")
	}

	// No id, users is been creating, so password is required.
	if u.Password == "" && u.UserID == "" {
		return NewMissingFieldError("password")
	}
	// Only check length if password is provided (creating or updating)
	if u.Password != "" && len(u.Password) < 6 {
		return &FieldError{
			UserError: UserError{
				Code:    "password_too_weak",
				Message: "Provided password need to be longer the 6 chars",
			},
			Field: "password",
		}
	}

	if u.Email.Valid && u.Email.String != "" && strings.Contains(u.Email.String, "@") {
		return &FieldError{
			UserError: UserError{
				Code:    "invalid_email",
				Message: "Provided email doesn't seems to be valid",
			},
			Field: "email",
		}
	}

	return nil
}

// ListOptions contains filtering, sorting and other fields to filter the list of user.
type ListOptions struct {
	Status  string
	Sort    string
	PerPage int64
	Page    int
	Cursor  string
}

func NewListOptions() *ListOptions {
	return &ListOptions{
		PerPage: 10,
	}
}

// ListResponse contains the list of users returned by List method.
type ListResponse struct {
	Total      int64   `json:"total"`
	PerPage    int64   `json:"per_page"`
	Users      []*User `json:"users"`
	NextCursor string  `json:"next_cursor"`
}
