package models

import "time"

type Expense struct {
	ID          int
	UserID      int
	Amount      float64
	Category    string
	Description string
	ExpenseDate time.Time
}
