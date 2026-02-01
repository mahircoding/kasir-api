package models

// Product represents a product in the store
type Product struct {
	ID         int     `json:"id"`
	Name       string  `json:"name"`
	Price      float64 `json:"price"`
	Stock      int     `json:"stock"`
	CategoryID int     `json:"category_id"`
}
