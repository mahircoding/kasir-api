package services

import (
	"kasir-api/models"
	"kasir-api/repositories"
)

// CategoryService handles business logic for categories
type CategoryService struct {
	repo *repositories.CategoryRepository
}

// NewCategoryService creates a new CategoryService
func NewCategoryService(repo *repositories.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

// GetAllCategories returns all categories
func (s *CategoryService) GetAllCategories() []models.Category {
	return s.repo.GetAll()
}

// GetCategoryByID returns a category by ID
func (s *CategoryService) GetCategoryByID(id int) (*models.Category, error) {
	return s.repo.GetByID(id)
}

// CreateCategory creates a new category
func (s *CategoryService) CreateCategory(category models.Category) models.Category {
	return s.repo.Create(category)
}

// UpdateCategory updates an existing category
func (s *CategoryService) UpdateCategory(id int, category models.Category) (*models.Category, error) {
	return s.repo.Update(id, category)
}

// DeleteCategory deletes a category by ID
func (s *CategoryService) DeleteCategory(id int) error {
	return s.repo.Delete(id)
}
