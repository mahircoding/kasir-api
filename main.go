package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"kasir-api/database"
	"kasir-api/handlers"
	"kasir-api/repositories"
	"kasir-api/services"

	"github.com/spf13/viper"

	_ "kasir-api/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

// AppConfig holds the application configuration
type AppConfig struct {
	Port   string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
}

// @title Kasir API
// @version 1.0
// @description API untuk sistem kasir sederhana
// @host localhost:8080
// @BasePath /api

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

func main() {
	// Load configuration from environment
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	port := viper.GetString("PORT")
	if port == "" {
		port = os.Getenv("PORT")
	}
	if port == "" {
		port = "8080"
	}

	// Check for DATABASE_URL first (Railway provides this)
	var db *sql.DB
	databaseURL := os.Getenv("DATABASE_URL")

	// Debug: print whether DATABASE_URL is set
	if databaseURL != "" {
		fmt.Println("DATABASE_URL is set, length:", len(databaseURL))
	} else {
		fmt.Println("DATABASE_URL is NOT set, checking individual vars...")
		fmt.Println("DB_HOST:", os.Getenv("DB_HOST"))
	}

	if databaseURL != "" {
		// Use DATABASE_URL directly (Railway/Supabase format)
		fmt.Println("Using DATABASE_URL connection string")
		var err error
		db, err = database.InitDB(databaseURL)
		if err != nil {
			log.Fatal("Failed to initialize database:", err)
		}
	} else {
		// Fallback to individual environment variables - use os.Getenv directly
		dbHost := os.Getenv("DB_HOST")
		dbPort := os.Getenv("DB_PORT")
		dbUser := os.Getenv("DB_USER")
		dbPassword := os.Getenv("DB_PASSWORD")
		dbName := os.Getenv("DB_NAME")
		dbSSLMode := os.Getenv("DB_SSLMODE")

		if dbPort == "" {
			dbPort = "5432"
		}
		if dbSSLMode == "" {
			dbSSLMode = "require"
		}

		dbConfig := database.DBConfig{
			Host:     dbHost,
			Port:     dbPort,
			User:     dbUser,
			Password: dbPassword,
			DBName:   dbName,
			SSLMode:  dbSSLMode,
		}

		var err error
		db, err = database.InitDBWithConfig(dbConfig)
		if err != nil {
			log.Fatal("Failed to initialize database:", err)
		}
	}
	defer db.Close()

	// Initialize product layers
	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	// Initialize category layers
	categoryRepo := repositories.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	// Define HTTP routes
	http.HandleFunc("/api/health", healthHandler)
	http.HandleFunc("/api/products", productHandler.Handle)
	http.HandleFunc("/api/products/", productHandler.Handle)
	http.HandleFunc("/api/categories", categoryHandler.Handle)
	http.HandleFunc("/api/categories/", categoryHandler.Handle)

	// Swagger documentation
	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	fmt.Println("Server running on localhost:" + port)
	fmt.Println("Swagger docs available at: http://localhost:" + port + "/swagger/index.html")
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

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
