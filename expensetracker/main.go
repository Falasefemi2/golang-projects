package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Expense struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Amount      int       `json:"amount"`
	Category    string    `json:"category,omitempty"`
	Date        time.Time `json:"date"`
}

type ExpenseStore struct {
	Expenses []Expense `json:"expenses"`
}

type MonthlyBudget struct {
	Month  int `json:"month"` // 1–12
	Year   int `json:"year"`
	Amount int `json:"amount"`
}

const (
	trackerFile = "tracker.json"
)

var (
	ErrEmptyDesc = errors.New("description cannot be empty")
	ErrBadAmount = errors.New("amount must be greater than zero")
	ErrNotFound  = errors.New("expense not found")
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	command := os.Args[1]

	switch command {
	case "add":
		if len(os.Args) < 4 {
			fmt.Println("Usage: expense-tracker add <description> <amount> [category]")
			return
		}

		description := os.Args[2]
		amount, err := strconv.Atoi(os.Args[3])
		if err != nil {
			fmt.Println("invalid amount")
		}
	}
}

func printUsage() {
	fmt.Println(`
Usage:
  expense-tracker add <description> <amount> [category]
      Add a new expense

  expense-tracker update <id> <description> <amount> [category]
      Update an expense

  expense-tracker delete <id>
      Delete an expense

  expense-tracker list
      View all expenses

  expense-tracker summary
      View total expenses

  expense-tracker summary <month>
      View expenses for a specific month (current year)
`)
}

func (es *ExpenseStore) getNextID() int {
	maxID := 0
	for _, e := range es.Expenses {
		if e.ID > maxID {
			maxID = e.ID
		}
	}
	return maxID + 1
}

func (es *ExpenseStore) AddExpense(description string, amount int, category string) error {
	if description == "" {
		return ErrEmptyDesc
	}
	if amount == 0 {
		return ErrBadAmount
	}
	id := es.getNextID()
	expense := Expense{
		ID:          id,
		Description: description,
		Amount:      amount,
		Category:    category,
		Date:        time.Now().UTC(),
	}
	es.Expenses = append(es.Expenses, expense)
	return nil
}

func (ex *ExpenseStore) DeleteExpense(id int) error {
	for i, expense := range ex.Expenses {
		if expense.ID == id {
			ex.Expenses = append(ex.Expenses[:i], ex.Expenses[i+1:]...)
			return nil
		}
	}
	return ErrNotFound
}

func (ex *ExpenseStore) ListExpense() {
	for _, expense := range ex.Expenses {
		fmt.Printf("ID: %d, Description: %s, Amount: %d, Category: %s, Date: %s\n",
			expense.ID, expense.Description, expense.Amount, expense.Category, expense.Date.Format("2006-01-02"))
	}
}

func (ex *ExpenseStore) UpdateExpense(id int, description string, amount int, category string) error {
	if description == "" {
		return ErrEmptyDesc
	}
	if amount == 0 {
		return ErrBadAmount
	}
	for i := range ex.Expenses {
		if ex.Expenses[i].ID == id {
			ex.Expenses[i].Description = description
			ex.Expenses[i].Amount = amount
			ex.Expenses[i].Category = category
			return nil
		}
	}
	return ErrNotFound
}

func loadExpenses() (*ExpenseStore, error) {
	data, err := os.ReadFile(trackerFile)
	if err != nil {
		if os.IsNotExist(err) {
			return &ExpenseStore{Expenses: []Expense{}}, nil
		}
		return nil, err
	}
	var expenseStore ExpenseStore
	if err := json.Unmarshal(data, &expenseStore); err != nil {
		return nil, err
	}
	return &expenseStore, nil
}

func saveExpenses(es *ExpenseStore) error {
	data, err := json.MarshalIndent(es, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(trackerFile, data, 0o644)
}
