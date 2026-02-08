package services

import (
	"kasir-api/models"
	"kasir-api/repositories"
	"time"
)

type ReportService struct {
	repo *repositories.ReportRepository
}

func NewReportService(repo *repositories.ReportRepository) *ReportService {
	return &ReportService{repo: repo}
}

func (s *ReportService) GetDailyReport(start, end string) (*models.SalesReport, error) {
	// Jika parameter kosong, set ke hari ini (YYYY-MM-DD)
	if start == "" || end == "" {
		today := time.Now().Format("2006-01-02")
		start = today + " 00:00:00"
		end = today + " 23:59:59"
	} else {
		start = start + " 00:00:00"
		end = end + " 23:59:59"
	}
	return s.repo.GetSalesReport(start, end)
}
