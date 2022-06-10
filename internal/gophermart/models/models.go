package models

import "time"

type User struct {
	ID       *int
	Login    string
	Password string
}

type Order struct {
	Number  string    `json:"number"`
	Status  string    `json:"status"`
	Accrual float32   `json:"accrual"`
	Upload  time.Time `json:"uploaded_at"`
}
