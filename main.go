package kasirapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Stock int     `json:"stock"`
}

var products = []Product{
	{ID: 1, Name: "Apple", Price: 0.5, Stock: 100},
	{ID: 2, Name: "Banana", Price: 0.3, Stock: 150},
	{ID: 3, Name: "Orange", Price: 0.7, Stock: 200},
}

func getProductByID(id int) (*Product, error) {
	for _, p := range products {
		if p.ID == id {
			return &p, nil
		}
	}
	return nil, fmt.Errorf("Product with ID %d not found", id)
}

func listProductsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func getProductHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := strings.TrimPrefix(r.URL.Path, "/products/")
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

func main() {
	http.HandleFunc("/products", listProductsHandler)
	http.HandleFunc("/products/", getProductHandler)

	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}		

