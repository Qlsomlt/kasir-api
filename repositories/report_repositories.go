package repositories

import (
	"database/sql"
	"kasir-api/models"
)

type ReportRepository struct {
	db *sql.DB
}

// FIX: Return harus *ReportRepository, dan return struct ReportRepository
func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

func (repo *ReportRepository) GetSalesReport(startDate, endDate string) (*models.SalesReport, error) {
	var report models.SalesReport

	// 1. Hitung Total Revenue & Total Transaksi
	querySummary := `
		SELECT COALESCE(SUM(total_amount), 0), COUNT(id) 
		FROM transactions 
		WHERE created_at BETWEEN $1 AND $2`

	err := repo.db.QueryRow(querySummary, startDate, endDate).Scan(&report.TotalRevenue, &report.TotalTransaksi)
	if err != nil {
		return nil, err
	}

	// 2. Cari Produk Terlaris
	queryBestSeller := `
		SELECT p.name, SUM(td.quantity) as total_qty
		FROM transaction_details td
		JOIN product p ON td.product_id = p.id  -- GANTI DI SINI (products -> product)
		JOIN transactions t ON td.transaction_id = t.id
		WHERE t.created_at BETWEEN $1 AND $2
		GROUP BY p.name
		ORDER BY total_qty DESC
		LIMIT 1`

	err = repo.db.QueryRow(queryBestSeller, startDate, endDate).Scan(
		&report.ProdukTerlaris.Nama,
		&report.ProdukTerlaris.QtyTerjual,
	)

	// Jika tidak ada transaksi sama sekali, set default agar tidak error
	if err == sql.ErrNoRows {
		report.ProdukTerlaris.Nama = "-"
		report.ProdukTerlaris.QtyTerjual = 0
	} else if err != nil {
		return nil, err
	}

	return &report, nil
}
