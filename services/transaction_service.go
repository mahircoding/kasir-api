package services

import (
	"fmt"

	"kasir-api/models"
	"kasir-api/repositories"
)

// TransactionService handles business logic for transactions
type TransactionService struct {
	transactionRepo *repositories.TransactionRepository
	productRepo     *repositories.ProductRepository
}

// NewTransactionService creates a new TransactionService
func NewTransactionService(transactionRepo *repositories.TransactionRepository, productRepo *repositories.ProductRepository) *TransactionService {
	return &TransactionService{
		transactionRepo: transactionRepo,
		productRepo:     productRepo,
	}
}

// CreateTransaction creates a new transaction from items
func (s *TransactionService) CreateTransaction(req models.CreateTransactionRequest) (*models.Transaction, error) {
	if len(req.Items) == 0 {
		return nil, fmt.Errorf("transaction must have at least one item")
	}

	var totalAmount int
	var details []models.TransactionDetail

	for _, item := range req.Items {
		// Get product to calculate subtotal
		product, err := s.productRepo.GetByID(item.ProductID)
		if err != nil {
			return nil, fmt.Errorf("product with ID %d not found", item.ProductID)
		}

		if item.Quantity <= 0 {
			return nil, fmt.Errorf("quantity must be greater than 0")
		}

		subtotal := int(product.Price) * item.Quantity
		totalAmount += subtotal

		details = append(details, models.TransactionDetail{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Subtotal:  subtotal,
		})
	}

	transaction := models.Transaction{
		TotalAmount: totalAmount,
		Details:     details,
	}

	return s.transactionRepo.Create(transaction)
}

// GetAllTransactions returns all transactions
func (s *TransactionService) GetAllTransactions() ([]models.Transaction, error) {
	return s.transactionRepo.GetAll()
}

// GetTransactionByID returns a transaction by ID
func (s *TransactionService) GetTransactionByID(id int) (*models.Transaction, error) {
	return s.transactionRepo.GetByID(id)
}

// DeleteTransaction deletes a transaction by ID
func (s *TransactionService) DeleteTransaction(id int) error {
	return s.transactionRepo.Delete(id)
}
