package model

import "time"

//Info from JWT token
type AuthUser struct {
	Username string
	Role     string
	Exp      time.Time
}

type User struct {
	Username string `json:"username"`
}

type Follows struct {
	From User `json:"from"`
	To   User `json:"to"`
}
