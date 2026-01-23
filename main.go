package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	_ "kasir-api/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Kasir API
// @version 1.0
// @description API untuk sistem kasir sederhana
// @host localhost:8080
// @BasePath /api

type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Product struct {
	ID         int     `json:"id"`
	Name       string  `json:"name"`
	Price      float64 `json:"price"`
	Stock      int     `json:"stock"`
	CategoryID int     `json:"category_id"`
}

var categories = []Category{
	{ID: 1, Name: "Buah", Description: "Kategori untuk berbagai jenis buah-buahan segar"},
	{ID: 2, Name: "Sayuran", Description: "Kategori untuk berbagai jenis sayuran segar"},
	{ID: 3, Name: "Minuman", Description: "Kategori untuk berbagai jenis minuman"},
}

var products = []Product{
	{ID: 1, Name: "Apple", Price: 90000, Stock: 100, CategoryID: 1},
	{ID: 2, Name: "Banana", Price: 30000, Stock: 150, CategoryID: 1},
	{ID: 3, Name: "Orange", Price: 70000, Stock: 200, CategoryID: 1},
}

func getProductByID(id int) (*Product, error) {
	for _, p := range products {
		if p.ID == id {
			return &p, nil
		}
	}
	return nil, fmt.Errorf("Product with ID %d not found", id)
}

func getProductIndex(id int) int {
	for i, p := range products {
		if p.ID == id {
			return i
		}
	}
	return -1
}

func generateProductID() int {
	maxID := 0
	for _, p := range products {
		if p.ID > maxID {
			maxID = p.ID
		}
	}
	return maxID + 1
}

func getCategoryByID(id int) (*Category, error) {
	for _, c := range categories {
		if c.ID == id {
			return &c, nil
		}
	}
	return nil, fmt.Errorf("Category with ID %d not found", id)
}

func getCategoryIndex(id int) int {
	for i, c := range categories {
		if c.ID == id {
			return i
		}
	}
	return -1
}

func generateCategoryID() int {
	maxID := 0
	for _, c := range categories {
		if c.ID > maxID {
			maxID = c.ID
		}
	}
	return maxID + 1
}

// healthHandler menampilkan status health API
// @Summary Health check
// @Description Check API health status
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "OK",
		"message": "API Running",
	})
}

// listProductsHandler menampilkan semua produk
// @Summary List all products
// @Description Get all products
// @Tags products
// @Produce json
// @Success 200 {array} Product
// @Router /products [get]
func listProductsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

// getProductHandler menampilkan detail produk berdasarkan ID
// @Summary Get product by ID
// @Description Get product details by ID
// @Tags products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} Product
// @Failure 400 {string} string "Invalid product ID"
// @Failure 404 {string} string "Product not found"
// @Router /products/{id} [get]
func getProductHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	product, err := getProductByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(product)
}

// createProductHandler membuat produk baru
// @Summary Create a new product
// @Description Create a new product
// @Tags products
// @Accept json
// @Produce json
// @Param product body Product true "Product object"
// @Success 201 {object} Product
// @Failure 400 {string} string "Invalid request body"
// @Router /products [post]
func createProductHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var newProduct Product
	err := json.NewDecoder(r.Body).Decode(&newProduct)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	newProduct.ID = generateProductID()
	products = append(products, newProduct)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newProduct)
}

// updateProductHandler mengupdate produk berdasarkan ID
// @Summary Update a product
// @Description Update product by ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param product body Product true "Product object"
// @Success 200 {object} Product
// @Failure 400 {string} string "Invalid request"
// @Failure 404 {string} string "Product not found"
// @Router /products/{id} [put]
func updateProductHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	index := getProductIndex(id)
	if index == -1 {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	var updatedProduct Product
	err = json.NewDecoder(r.Body).Decode(&updatedProduct)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	updatedProduct.ID = id
	products[index] = updatedProduct

	json.NewEncoder(w).Encode(updatedProduct)
}

// deleteProductHandler menghapus produk berdasarkan ID
// @Summary Delete a product
// @Description Delete product by ID
// @Tags products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {string} string "Product deleted successfully"
// @Failure 400 {string} string "Invalid product ID"
// @Failure 404 {string} string "Product not found"
// @Router /products/{id} [delete]
func deleteProductHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	index := getProductIndex(id)
	if index == -1 {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	products = append(products[:index], products[index+1:]...)

	json.NewEncoder(w).Encode(map[string]string{"message": "Product deleted successfully"})
}

// listCategoriesHandler menampilkan semua kategori
// @Summary List all categories
// @Description Get all categories
// @Tags categories
// @Produce json
// @Success 200 {array} Category
// @Router /categories [get]
func listCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

// getCategoryHandler menampilkan detail kategori berdasarkan ID
// @Summary Get category by ID
// @Description Get category details by ID
// @Tags categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} Category
// @Failure 400 {string} string "Invalid category ID"
// @Failure 404 {string} string "Category not found"
// @Router /categories/{id} [get]
func getCategoryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	category, err := getCategoryByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(category)
}

// createCategoryHandler membuat kategori baru
// @Summary Create a new category
// @Description Create a new category
// @Tags categories
// @Accept json
// @Produce json
// @Param category body Category true "Category object"
// @Success 201 {object} Category
// @Failure 400 {string} string "Invalid request body"
// @Router /categories [post]
func createCategoryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var newCategory Category
	err := json.NewDecoder(r.Body).Decode(&newCategory)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	newCategory.ID = generateCategoryID()
	categories = append(categories, newCategory)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newCategory)
}

// updateCategoryHandler mengupdate kategori berdasarkan ID
// @Summary Update a category
// @Description Update category by ID
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param category body Category true "Category object"
// @Success 200 {object} Category
// @Failure 400 {string} string "Invalid request"
// @Failure 404 {string} string "Category not found"
// @Router /categories/{id} [put]
func updateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	index := getCategoryIndex(id)
	if index == -1 {
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}

	var updatedCategory Category
	err = json.NewDecoder(r.Body).Decode(&updatedCategory)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	updatedCategory.ID = id
	categories[index] = updatedCategory

	json.NewEncoder(w).Encode(updatedCategory)
}

// deleteCategoryHandler menghapus kategori berdasarkan ID
// @Summary Delete a category
// @Description Delete category by ID
// @Tags categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {string} string "Category deleted successfully"
// @Failure 400 {string} string "Invalid category ID"
// @Failure 404 {string} string "Category not found"
// @Router /categories/{id} [delete]
func deleteCategoryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	index := getCategoryIndex(id)
	if index == -1 {
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}

	categories = append(categories[:index], categories[index+1:]...)

	json.NewEncoder(w).Encode(map[string]string{"message": "Category deleted successfully"})
}

// categoryHandler menangani routing berdasarkan method HTTP
func categoryHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if strings.TrimPrefix(r.URL.Path, "/api/categories") == "" || strings.TrimPrefix(r.URL.Path, "/api/categories") == "/" {
			listCategoriesHandler(w, r)
		} else {
			getCategoryHandler(w, r)
		}
	case http.MethodPost:
		createCategoryHandler(w, r)
	case http.MethodPut:
		updateCategoryHandler(w, r)
	case http.MethodDelete:
		deleteCategoryHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// productHandler menangani routing berdasarkan method HTTP
func productHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if strings.TrimPrefix(r.URL.Path, "/api/products") == "" || strings.TrimPrefix(r.URL.Path, "/api/products") == "/" {
			listProductsHandler(w, r)
		} else {
			getProductHandler(w, r)
		}
	case http.MethodPost:
		createProductHandler(w, r)
	case http.MethodPut:
		updateProductHandler(w, r)
	case http.MethodDelete:
		deleteProductHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	// Define HTTP routes
	http.HandleFunc("/api/health", healthHandler)
	http.HandleFunc("/api/products", productHandler)
	http.HandleFunc("/api/products/", productHandler)
	http.HandleFunc("/api/categories", categoryHandler)
	http.HandleFunc("/api/categories/", categoryHandler)

	// Swagger documentation
	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	fmt.Println("Server is running on port 8080")
	fmt.Println("Swagger docs available at: http://localhost:8080/swagger/index.html")
	fmt.Println("\nAvailable endpoints:")
	fmt.Println("Health:")
	fmt.Println("  GET    /api/health       - API health check")
	fmt.Println("\nProducts:")
	fmt.Println("  GET    /api/products     - List all products")
	fmt.Println("  GET    /api/products/{id} - Get product by ID")
	fmt.Println("  POST   /api/products     - Create new product")
	fmt.Println("  PUT    /api/products/{id} - Update product")
	fmt.Println("  DELETE /api/products/{id} - Delete product")
	fmt.Println("\nCategories:")
	fmt.Println("  GET    /api/categories     - List all categories")
	fmt.Println("  GET    /api/categories/{id} - Get category by ID")
	fmt.Println("  POST   /api/categories     - Create new category")
	fmt.Println("  PUT    /api/categories/{id} - Update category")
	fmt.Println("  DELETE /api/categories/{id} - Delete category")
	http.ListenAndServe(":8080", nil)
}
