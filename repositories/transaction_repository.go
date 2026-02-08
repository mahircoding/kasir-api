package repositories

import (
	"database/sql"
	"fmt"

	"kasir-api/models"
)

// TransactionRepository handles data access for transactions
type TransactionRepository struct {
	db *sql.DB
}

// NewTransactionRepository creates a new TransactionRepository
func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

// Create creates a new transaction with details
func (r *TransactionRepository) Create(transaction models.Transaction) (*models.Transaction, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Insert transaction
	err = tx.QueryRow(
		"INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id, created_at",
		transaction.TotalAmount,
	).Scan(&transaction.ID, &transaction.CreatedAt)
	if err != nil {
		return nil, err
	}

	// Prepare statement for inserting transaction details (more efficient for multiple inserts)
	stmt, err := tx.Prepare(`
		INSERT INTO transaction_details 
		(transaction_id, product_id, quantity, subtotal)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Insert transaction details using prepared statement
	for i := range transaction.Details {
		transaction.Details[i].TransactionID = transaction.ID
		err = stmt.QueryRow(
			transaction.ID,
			transaction.Details[i].ProductID,
			transaction.Details[i].Quantity,
			transaction.Details[i].Subtotal,
		).Scan(&transaction.Details[i].ID)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &transaction, nil
}

// GetAll returns all transactions
func (r *TransactionRepository) GetAll() ([]models.Transaction, error) {
	rows, err := r.db.Query("SELECT id, total_amount, created_at FROM transactions ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var t models.Transaction
		if err := rows.Scan(&t.ID, &t.TotalAmount, &t.CreatedAt); err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}
	return transactions, nil
}

// GetByID returns a transaction by ID with its details
func (r *TransactionRepository) GetByID(id int) (*models.Transaction, error) {
	var t models.Transaction
	err := r.db.QueryRow(
		"SELECT id, total_amount, created_at FROM transactions WHERE id = $1",
		id,
	).Scan(&t.ID, &t.TotalAmount, &t.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Transaction with ID %d not found", id)
		}
		return nil, err
	}

	// Get transaction details
	rows, err := r.db.Query(
		"SELECT id, transaction_id, product_id, quantity, subtotal FROM transaction_details WHERE transaction_id = $1",
		id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var d models.TransactionDetail
		if err := rows.Scan(&d.ID, &d.TransactionID, &d.ProductID, &d.Quantity, &d.Subtotal); err != nil {
			return nil, err
		}
		t.Details = append(t.Details, d)
	}

	return &t, nil
}

// Delete deletes a transaction by ID
func (r *TransactionRepository) Delete(id int) error {
	result, err := r.db.Exec("DELETE FROM transactions WHERE id = $1", id)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("Transaction with ID %d not found", id)
	}
	return nil
}
