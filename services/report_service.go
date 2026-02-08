package services

import (
	"time"

	"kasir-api/models"
	"kasir-api/repositories"
)

// ReportService handles business logic for reports
type ReportService struct {
	repo *repositories.ReportRepository
}

// NewReportService creates a new ReportService
func NewReportService(repo *repositories.ReportRepository) *ReportService {
	return &ReportService{repo: repo}
}

// GetTodayReport returns sales summary for today
func (s *ReportService) GetTodayReport() (*models.SalesReport, error) {
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	report, err := s.repo.GetSalesReport(startOfDay, endOfDay)
	if err != nil {
		return nil, err
	}

	report.StartDate = startOfDay.Format("2006-01-02")
	report.EndDate = now.Format("2006-01-02")

	return report, nil
}

// GetReportByDateRange returns sales summary for a date range
func (s *ReportService) GetReportByDateRange(startDateStr, endDateStr string) (*models.SalesReport, error) {
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return nil, err
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return nil, err
	}

	// Add 1 day to end date to include the entire end day
	endDate = endDate.Add(24 * time.Hour)

	report, err := s.repo.GetSalesReport(startDate, endDate)
	if err != nil {
		return nil, err
	}

	report.StartDate = startDateStr
	report.EndDate = endDateStr

	return report, nil
}
