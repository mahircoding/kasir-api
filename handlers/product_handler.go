package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"kasir-api/models"
	"kasir-api/services"
)

// ProductHandler handles HTTP requests for products
type ProductHandler struct {
	service *services.ProductService
}

// NewProductHandler creates a new ProductHandler
func NewProductHandler(service *services.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

// Handle menangani routing berdasarkan method HTTP
func (h *ProductHandler) Handle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		path := strings.TrimPrefix(r.URL.Path, "/api/products")
		if path == "" || path == "/" {
			h.ListProducts(w, r)
		} else {
			h.GetProduct(w, r)
		}
	case http.MethodPost:
		h.CreateProduct(w, r)
	case http.MethodPut:
		h.UpdateProduct(w, r)
	case http.MethodDelete:
		h.DeleteProduct(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// ListProducts menampilkan semua produk
// @Summary List all products
// @Description Get all products
// @Tags products
// @Produce json
// @Success 200 {array} models.Product
// @Router /products [get]
func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	products, err := h.service.GetAllProducts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(products)
}

// GetProduct menampilkan detail produk berdasarkan ID
// @Summary Get product by ID
// @Description Get product details by ID
// @Tags products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} models.Product
// @Failure 400 {string} string "Invalid product ID"
// @Failure 404 {string} string "Product not found"
// @Router /products/{id} [get]
func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	product, err := h.service.GetProductByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(product)
}

// CreateProduct membuat produk baru
// @Summary Create a new product
// @Description Create a new product
// @Tags products
// @Accept json
// @Produce json
// @Param product body models.Product true "Product object"
// @Success 201 {object} models.Product
// @Failure 400 {string} string "Invalid request body"
// @Router /products [post]
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var newProduct models.Product
	err := json.NewDecoder(r.Body).Decode(&newProduct)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	createdProduct, err := h.service.CreateProduct(newProduct)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdProduct)
}

// UpdateProduct mengupdate produk berdasarkan ID
// @Summary Update a product
// @Description Update product by ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param product body models.Product true "Product object"
// @Success 200 {object} models.Product
// @Failure 400 {string} string "Invalid request"
// @Failure 404 {string} string "Product not found"
// @Router /products/{id} [put]
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	var updatedProduct models.Product
	err = json.NewDecoder(r.Body).Decode(&updatedProduct)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	product, err := h.service.UpdateProduct(id, updatedProduct)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(product)
}

// DeleteProduct menghapus produk berdasarkan ID
// @Summary Delete a product
// @Description Delete product by ID
// @Tags products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {string} string "Product deleted successfully"
// @Failure 400 {string} string "Invalid product ID"
// @Failure 404 {string} string "Product not found"
// @Router /products/{id} [delete]
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteProduct(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Product deleted successfully"})
}
