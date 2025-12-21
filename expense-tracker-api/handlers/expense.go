package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"expense-tracker/internal/db"
	"expense-tracker/middleware"
	"expense-tracker/utils"
)

type expenseReq struct {
	Amount      float64 `json:"amount"`
	Category    string  `json:"category"`
	Description string  `json:"description"`
	ExpenseDate string  `json:"expense_date"`
}

func CreateExpense(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)

	var req expenseReq
	json.NewDecoder(r.Body).Decode(&req)

	date, _ := time.Parse("2006-01-02", req.ExpenseDate)

	db.DB.Exec(
		`INSERT INTO expenses (user_id, amount, category, description, expense_date)
		 VALUES (?, ?, ?, ?, ?)`,
		userID, req.Amount, req.Category, req.Description, date,
	)

	utils.JSON(w, http.StatusCreated, map[string]string{"message": "expense added"})
}

func ListExpenses(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)

	rows, _ := db.DB.Query(
		`SELECT id, amount, category, description, expense_date
		 FROM expenses WHERE user_id = ?`,
		userID,
	)

	var expenses []map[string]interface{}
	for rows.Next() {
		var id int
		var amount float64
		var category, desc string
		var date string

		rows.Scan(&id, &amount, &category, &desc, &date)

		expenses = append(expenses, map[string]interface{}{
			"id": id, "amount": amount, "category": category,
			"description": desc, "expense_date": date,
		})
	}

	utils.JSON(w, http.StatusOK, expenses)
}
