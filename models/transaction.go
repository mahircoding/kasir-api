package models

import "time"

// Transaction represents a sales transaction
type Transaction struct {
	ID          int                 `json:"id"`
	TotalAmount int                 `json:"total_amount"`
	CreatedAt   time.Time           `json:"created_at"`
	Details     []TransactionDetail `json:"details,omitempty"`
}

// TransactionDetail represents a detail line item in a transaction
type TransactionDetail struct {
	ID            int `json:"id"`
	TransactionID int `json:"transaction_id"`
	ProductID     int `json:"product_id"`
	Quantity      int `json:"quantity"`
	Subtotal      int `json:"subtotal"`
}

// CreateTransactionRequest represents the request body for creating a transaction
type CreateTransactionRequest struct {
	Items []TransactionItem `json:"items"`
}

// TransactionItem represents a single item in a transaction request
type TransactionItem struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}
