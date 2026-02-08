package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"kasir-api/services"
)

// ReportHandler handles HTTP requests for reports
type ReportHandler struct {
	service *services.ReportService
}

// NewReportHandler creates a new ReportHandler
func NewReportHandler(service *services.ReportService) *ReportHandler {
	return &ReportHandler{service: service}
}

// Handle menangani routing berdasarkan path
func (h *ReportHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/report")

	switch {
	case path == "/hari-ini" || path == "/hari-ini/":
		h.GetTodayReport(w, r)
	case path == "" || path == "/":
		h.GetReportByDateRange(w, r)
	default:
		http.Error(w, "Not found", http.StatusNotFound)
	}
}

// GetTodayReport menampilkan laporan penjualan hari ini
// @Summary Get today's sales report
// @Description Get sales summary for today including total revenue, transaction count, and best selling product
// @Tags report
// @Produce json
// @Success 200 {object} models.SalesReport
// @Router /report/hari-ini [get]
func (h *ReportHandler) GetTodayReport(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	report, err := h.service.GetTodayReport()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(report)
}

// GetReportByDateRange menampilkan laporan penjualan berdasarkan rentang tanggal
// @Summary Get sales report by date range
// @Description Get sales summary for a specific date range
// @Tags report
// @Produce json
// @Param start_date query string true "Start date (YYYY-MM-DD)"
// @Param end_date query string true "End date (YYYY-MM-DD)"
// @Success 200 {object} models.SalesReport
// @Failure 400 {string} string "Missing or invalid date parameters"
// @Router /report [get]
func (h *ReportHandler) GetReportByDateRange(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	if startDate == "" || endDate == "" {
		http.Error(w, "start_date and end_date are required (format: YYYY-MM-DD)", http.StatusBadRequest)
		return
	}

	report, err := h.service.GetReportByDateRange(startDate, endDate)
	if err != nil {
		http.Error(w, "Invalid date format. Use YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(report)
}
