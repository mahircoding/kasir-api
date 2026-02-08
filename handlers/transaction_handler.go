package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"kasir-api/models"
	"kasir-api/services"
)

// TransactionHandler handles HTTP requests for transactions
type TransactionHandler struct {
	service *services.TransactionService
}

// NewTransactionHandler creates a new TransactionHandler
func NewTransactionHandler(service *services.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: service}
}

// Handle menangani routing berdasarkan method HTTP
func (h *TransactionHandler) Handle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		path := strings.TrimPrefix(r.URL.Path, "/api/transactions")
		if path == "" || path == "/" {
			h.ListTransactions(w, r)
		} else {
			h.GetTransaction(w, r)
		}
	case http.MethodPost:
		h.CreateTransaction(w, r)
	case http.MethodDelete:
		h.DeleteTransaction(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// ListTransactions menampilkan semua transaksi
// @Summary List all transactions
// @Description Get all transactions
// @Tags transactions
// @Produce json
// @Success 200 {array} models.Transaction
// @Router /transactions [get]
func (h *TransactionHandler) ListTransactions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	transactions, err := h.service.GetAllTransactions()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(transactions)
}

// GetTransaction menampilkan detail transaksi berdasarkan ID
// @Summary Get transaction by ID
// @Description Get transaction details by ID including line items
// @Tags transactions
// @Produce json
// @Param id path int true "Transaction ID"
// @Success 200 {object} models.Transaction
// @Failure 400 {string} string "Invalid transaction ID"
// @Failure 404 {string} string "Transaction not found"
// @Router /transactions/{id} [get]
func (h *TransactionHandler) GetTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := strings.TrimPrefix(r.URL.Path, "/api/transactions/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid transaction ID", http.StatusBadRequest)
		return
	}

	transaction, err := h.service.GetTransactionByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(transaction)
}

// CreateTransaction membuat transaksi baru
// @Summary Create a new transaction
// @Description Create a new transaction with items
// @Tags transactions
// @Accept json
// @Produce json
// @Param transaction body models.CreateTransactionRequest true "Transaction items"
// @Success 201 {object} models.Transaction
// @Failure 400 {string} string "Invalid request body"
// @Router /transactions [post]
func (h *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req models.CreateTransactionRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	transaction, err := h.service.CreateTransaction(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(transaction)
}

// DeleteTransaction menghapus transaksi berdasarkan ID
// @Summary Delete a transaction
// @Description Delete transaction by ID
// @Tags transactions
// @Produce json
// @Param id path int true "Transaction ID"
// @Success 200 {string} string "Transaction deleted successfully"
// @Failure 400 {string} string "Invalid transaction ID"
// @Failure 404 {string} string "Transaction not found"
// @Router /transactions/{id} [delete]
func (h *TransactionHandler) DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := strings.TrimPrefix(r.URL.Path, "/api/transactions/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid transaction ID", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteTransaction(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Transaction deleted successfully"})
}
