package models

import (
	"time"
)

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	APIKey   string
}

type OrderType string

const (
	New        string = "NEW"
	Processing string = "PROCESSING"
)

type Order struct {
	Number    string    `json:"number"`
	Status    string    `json:"status"`
	Accrual   float32   `json:"accrual,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type Balance struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

type Withdraw struct {
	Order       string    `json:"order"`
	Sum         float64   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at,omitempty"`
}
