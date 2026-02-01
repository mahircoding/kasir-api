package repositories

import (
	"fmt"

	"kasir-api/models"
)

// CategoryRepository handles data access for categories
type CategoryRepository struct {
	categories []models.Category
}

// NewCategoryRepository creates a new CategoryRepository with initial data
func NewCategoryRepository() *CategoryRepository {
	return &CategoryRepository{
		categories: []models.Category{
			{ID: 1, Name: "Buah", Description: "Kategori untuk berbagai jenis buah-buahan segar"},
			{ID: 2, Name: "Sayuran", Description: "Kategori untuk berbagai jenis sayuran segar"},
			{ID: 3, Name: "Minuman", Description: "Kategori untuk berbagai jenis minuman"},
		},
	}
}

// GetAll returns all categories
func (r *CategoryRepository) GetAll() []models.Category {
	return r.categories
}

// GetByID returns a category by ID
func (r *CategoryRepository) GetByID(id int) (*models.Category, error) {
	for _, c := range r.categories {
		if c.ID == id {
			return &c, nil
		}
	}
	return nil, fmt.Errorf("Category with ID %d not found", id)
}

// GetIndex returns the index of a category by ID
func (r *CategoryRepository) GetIndex(id int) int {
	for i, c := range r.categories {
		if c.ID == id {
			return i
		}
	}
	return -1
}

// GenerateID generates a new unique ID for a category
func (r *CategoryRepository) GenerateID() int {
	maxID := 0
	for _, c := range r.categories {
		if c.ID > maxID {
			maxID = c.ID
		}
	}
	return maxID + 1
}

// Create adds a new category
func (r *CategoryRepository) Create(category models.Category) models.Category {
	category.ID = r.GenerateID()
	r.categories = append(r.categories, category)
	return category
}

// Update updates an existing category
func (r *CategoryRepository) Update(id int, category models.Category) (*models.Category, error) {
	index := r.GetIndex(id)
	if index == -1 {
		return nil, fmt.Errorf("Category with ID %d not found", id)
	}
	category.ID = id
	r.categories[index] = category
	return &category, nil
}

// Delete removes a category by ID
func (r *CategoryRepository) Delete(id int) error {
	index := r.GetIndex(id)
	if index == -1 {
		return fmt.Errorf("Category with ID %d not found", id)
	}
	r.categories = append(r.categories[:index], r.categories[index+1:]...)
	return nil
}
