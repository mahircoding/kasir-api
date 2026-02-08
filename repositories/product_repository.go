package repositories

import (
	"database/sql"
	"fmt"

	"kasir-api/models"
)

// ProductRepository handles data access for products
type ProductRepository struct {
	db *sql.DB
}

// NewProductRepository creates a new ProductRepository
func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// GetAll returns all products with optional filters
func (r *ProductRepository) GetAll(filter models.ProductFilter) ([]models.Product, error) {
	query := "SELECT id, name, price, stock, category_id FROM products WHERE 1=1"
	var args []interface{}
	argIndex := 1

	if filter.Name != "" {
		query += fmt.Sprintf(" AND LOWER(name) LIKE LOWER($%d)", argIndex)
		args = append(args, "%"+filter.Name+"%")
		argIndex++
	}

	if filter.CategoryID > 0 {
		query += fmt.Sprintf(" AND category_id = $%d", argIndex)
		args = append(args, filter.CategoryID)
		argIndex++
	}

	if filter.MinPrice > 0 {
		query += fmt.Sprintf(" AND price >= $%d", argIndex)
		args = append(args, filter.MinPrice)
		argIndex++
	}

	if filter.MaxPrice > 0 {
		query += fmt.Sprintf(" AND price <= $%d", argIndex)
		args = append(args, filter.MaxPrice)
		argIndex++
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.CategoryID); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

// GetByID returns a product by ID
func (r *ProductRepository) GetByID(id int) (*models.Product, error) {
	var p models.Product
	err := r.db.QueryRow("SELECT id, name, price, stock, category_id FROM products WHERE id = $1", id).
		Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.CategoryID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Product with ID %d not found", id)
		}
		return nil, err
	}
	return &p, nil
}

// Create adds a new product
func (r *ProductRepository) Create(product models.Product) (*models.Product, error) {
	err := r.db.QueryRow(
		"INSERT INTO products (name, price, stock, category_id) VALUES ($1, $2, $3, $4) RETURNING id",
		product.Name, product.Price, product.Stock, product.CategoryID,
	).Scan(&product.ID)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// Update updates an existing product
func (r *ProductRepository) Update(id int, product models.Product) (*models.Product, error) {
	result, err := r.db.Exec(
		"UPDATE products SET name = $1, price = $2, stock = $3, category_id = $4 WHERE id = $5",
		product.Name, product.Price, product.Stock, product.CategoryID, id,
	)
	if err != nil {
		return nil, err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return nil, fmt.Errorf("Product with ID %d not found", id)
	}
	product.ID = id
	return &product, nil
}

// Delete removes a product by ID
func (r *ProductRepository) Delete(id int) error {
	result, err := r.db.Exec("DELETE FROM products WHERE id = $1", id)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("Product with ID %d not found", id)
	}
	return nil
}
