package entities

import "time"

type Order struct {
	ID        string    `json:"id"`
	Product   string    `json:"product"`
	Price     float64   `json:"price"`
	Priority  int       `json:"priority"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}
