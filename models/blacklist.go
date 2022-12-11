package models

import "time"

// Blacklist is a root struct that is used to store the json encoded data for/from a mongodb blacklist doc.
type Blacklist struct {
	Id        string    `json:"id,omitempty"`
	UserID    string    `json:"user_id,omitempty"`
	AuthToken string    `json:"auth_token,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}
