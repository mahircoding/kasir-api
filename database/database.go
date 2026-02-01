package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// DB holds the database connection
var DB *sql.DB

// DBConfig holds database configuration
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// BuildConnectionString builds a PostgreSQL connection string from config
func BuildConnectionString(cfg DBConfig) string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)
}

// InitDB initializes the database connection
func InitDB(connStr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	DB = db
	fmt.Println("Database connected successfully")
	return db, nil
}

// InitDBWithConfig initializes the database connection using DBConfig
func InitDBWithConfig(cfg DBConfig) (*sql.DB, error) {
	connStr := BuildConnectionString(cfg)
	return InitDB(connStr)
}
