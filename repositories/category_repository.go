package repositories

import (
	"database/sql"
	"fmt"

	"kasir-api/models"
)

// CategoryRepository handles data access for categories
type CategoryRepository struct {
	db *sql.DB
}

// NewCategoryRepository creates a new CategoryRepository
func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

// GetAll returns all categories
func (r *CategoryRepository) GetAll() ([]models.Category, error) {
	rows, err := r.db.Query("SELECT id, name, description FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var c models.Category
		var description sql.NullString
		if err := rows.Scan(&c.ID, &c.Name, &description); err != nil {
			return nil, err
		}
		if description.Valid {
			c.Description = description.String
		}
		categories = append(categories, c)
	}
	return categories, nil
}

// GetByID returns a category by ID
func (r *CategoryRepository) GetByID(id int) (*models.Category, error) {
	var c models.Category
	var description sql.NullString
	err := r.db.QueryRow("SELECT id, name, description FROM categories WHERE id = $1", id).
		Scan(&c.ID, &c.Name, &description)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Category with ID %d not found", id)
		}
		return nil, err
	}
	if description.Valid {
		c.Description = description.String
	}
	return &c, nil
}

// Create adds a new category
func (r *CategoryRepository) Create(category models.Category) (*models.Category, error) {
	err := r.db.QueryRow(
		"INSERT INTO categories (name, description) VALUES ($1, $2) RETURNING id",
		category.Name, category.Description,
	).Scan(&category.ID)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// Update updates an existing category
func (r *CategoryRepository) Update(id int, category models.Category) (*models.Category, error) {
	result, err := r.db.Exec(
		"UPDATE categories SET name = $1, description = $2 WHERE id = $3",
		category.Name, category.Description, id,
	)
	if err != nil {
		return nil, err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return nil, fmt.Errorf("Category with ID %d not found", id)
	}
	category.ID = id
	return &category, nil
}

// Delete removes a category by ID
func (r *CategoryRepository) Delete(id int) error {
	result, err := r.db.Exec("DELETE FROM categories WHERE id = $1", id)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("Category with ID %d not found", id)
	}
	return nil
}
