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

type Balance struct {
	Balance   float32 `json:"current"`
	Withdraws float32 `json:"withdrawn"`
}

type Withdraw struct {
	Number    string    `json:"number"`
	Processed time.Time `json:"processed_at"`
	Withdraw  float32   `json:"sum"`
}

type Accrual struct {
	Order  string  `json:"order"`
	Status string  `json:"status"`
	Value  float32 `json:"accrual,omitempty"`
}

type OrderToProcess struct {
	Login string
	Order string
}
