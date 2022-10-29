package user

import "time"

// User object.
type User struct {
	Id           string    `json:"id"`
	FirstName    *string   `json:"firstName"`
	LastName     *string   `json:"lastName"`
	Email        string    `json:"email"`
	Role         Role      `json:"type"`
	PasswordHash string    `json:"_"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}
