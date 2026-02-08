package services

import (
	"kasir-api/models"
	"kasir-api/repositories"
)

// ProductService handles business logic for products
type ProductService struct {
	repo *repositories.ProductRepository
}

// NewProductService creates a new ProductService
func NewProductService(repo *repositories.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

// GetAllProducts returns all products with optional filters
func (s *ProductService) GetAllProducts(filter models.ProductFilter) ([]models.Product, error) {
	return s.repo.GetAll(filter)
}

// GetProductByID returns a product by ID
func (s *ProductService) GetProductByID(id int) (*models.Product, error) {
	return s.repo.GetByID(id)
}

// CreateProduct creates a new product
func (s *ProductService) CreateProduct(product models.Product) (*models.Product, error) {
	return s.repo.Create(product)
}

// UpdateProduct updates an existing product
func (s *ProductService) UpdateProduct(id int, product models.Product) (*models.Product, error) {
	return s.repo.Update(id, product)
}

// DeleteProduct deletes a product by ID
func (s *ProductService) DeleteProduct(id int) error {
	return s.repo.Delete(id)
}
