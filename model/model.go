package model

import "time"

// Info from JWT token
type AuthUser struct {
	Username string
	Role     string
	Exp      time.Time
}

type User struct {
	Username    string `json:"username"`
	Town        string `json:"town"`
	Gender      string `json:"gender"`
	YearOfBirth int32  `json:"yearOfBirth"`
	IsPrivate   bool   `json:"private"`
}

type Approved struct {
	Approved bool `json:"approved"`
}

type TargetUserGroup struct {
	Town   string `json:"town"`
	Gender string `json:"gender"`
	MinAge int32  `json:"min-age"`
	MaxAge int32  `json:"max-age"`
}
