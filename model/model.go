package model

type User struct {
	Username string `json:"username"`
}

type Follows struct {
	From User `json:"from"`
	To   User `json:"to"`
}

type Approved struct {
	Approved bool `json:"approved"`
}
