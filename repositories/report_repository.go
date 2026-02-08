package repositories

import (
	"database/sql"
	"time"

	"kasir-api/models"
)

// ReportRepository handles data access for reports
type ReportRepository struct {
	db *sql.DB
}

// NewReportRepository creates a new ReportRepository
func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

// GetSalesReport returns sales summary for a date range
func (r *ReportRepository) GetSalesReport(startDate, endDate time.Time) (*models.SalesReport, error) {
	report := &models.SalesReport{}

	// Get total revenue and transaction count
	err := r.db.QueryRow(`
		SELECT 
			COALESCE(SUM(total_amount), 0) as total_revenue,
			COUNT(*) as total_transaksi
		FROM transactions
		WHERE created_at >= $1 AND created_at < $2
	`, startDate, endDate).Scan(&report.TotalRevenue, &report.TotalTransaksi)
	if err != nil {
		return nil, err
	}

	// Get best selling product
	var productName sql.NullString
	var qtyTerjual sql.NullInt64
	err = r.db.QueryRow(`
		SELECT 
			p.name,
			SUM(td.quantity) as qty_terjual
		FROM transaction_details td
		JOIN transactions t ON td.transaction_id = t.id
		JOIN products p ON td.product_id = p.id
		WHERE t.created_at >= $1 AND t.created_at < $2
		GROUP BY p.id, p.name
		ORDER BY qty_terjual DESC
		LIMIT 1
	`, startDate, endDate).Scan(&productName, &qtyTerjual)

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if productName.Valid && qtyTerjual.Valid {
		report.ProdukTerlaris = &models.BestSellerInfo{
			Nama:       productName.String,
			QtyTerjual: int(qtyTerjual.Int64),
		}
	}

	return report, nil
}
