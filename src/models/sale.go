package models

type Sale struct {
	ID              string `json:"id"`
	UserId          string `json:"user_id"`
	StoreId         string `json:"store_id"`
	Date            string `json:"date"`
	SalesInCents    int    `json:"sales_in_cents"`
	ExpensesInCents int    `json:"expenses_in_cents"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}
