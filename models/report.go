package models

// SalesReport represents the sales summary report
type SalesReport struct {
	TotalRevenue   int             `json:"total_revenue"`
	TotalTransaksi int             `json:"total_transaksi"`
	ProdukTerlaris *BestSellerInfo `json:"produk_terlaris,omitempty"`
	StartDate      string          `json:"start_date,omitempty"`
	EndDate        string          `json:"end_date,omitempty"`
}

// BestSellerInfo represents the best selling product info
type BestSellerInfo struct {
	Nama       string `json:"nama"`
	QtyTerjual int    `json:"qty_terjual"`
}
