package model

import "time"

// Info from JWT token
type AuthUser struct {
	Username string
	Role     string
	Exp      time.Time
}

type User struct {
	Username  string `json:"username"`
	IsPrivate bool   `json:"private"`
}

type Approved struct {
	Approved bool `json:"approved"`
}
