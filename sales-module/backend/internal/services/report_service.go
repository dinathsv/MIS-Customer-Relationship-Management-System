package services

import (
	"database/sql"
	"fmt"
	"sales-module/internal/models"
)

type ReportService struct {
	db *sql.DB
}

func NewReportService(db *sql.DB) *ReportService {
	return &ReportService{db: db}
}

func (s *ReportService) GetDashboardStats() (*models.DashboardStats, error) {
	stats := &models.DashboardStats{}

	s.db.QueryRow(`SELECT COALESCE(SUM(total_amount),0), COUNT(*) FROM sales WHERE status='completed'`).Scan(&stats.TotalRevenue, &stats.TotalSales)
	s.db.QueryRow(`SELECT COUNT(*) FROM sales WHERE status='pending'`).Scan(&stats.PendingOrders)

	if stats.TotalSales > 0 {
		stats.AverageOrderValue = stats.TotalRevenue / float64(stats.TotalSales)
	}

	// Growth: compare last 30 days vs previous 30 days
	var currentRev, prevRev float64
	var currentCount, prevCount int
	s.db.QueryRow(`SELECT COALESCE(SUM(total_amount),0), COUNT(*) FROM sales WHERE status='completed' AND created_at >= NOW()-INTERVAL '30 days'`).Scan(&currentRev, &currentCount)
	s.db.QueryRow(`SELECT COALESCE(SUM(total_amount),0), COUNT(*) FROM sales WHERE status='completed' AND created_at >= NOW()-INTERVAL '60 days' AND created_at < NOW()-INTERVAL '30 days'`).Scan(&prevRev, &prevCount)

	if prevRev > 0 { stats.RevenueGrowth = ((currentRev - prevRev) / prevRev) * 100 }
	if prevCount > 0 { stats.SalesGrowth = ((float64(currentCount) - float64(prevCount)) / float64(prevCount)) * 100 }

	return stats, nil
}

func (s *ReportService) GetSalesSummary(period string) ([]models.SalesSummary, error) {
	var dateFormat, groupBy string
	switch period {
	case "daily":
		dateFormat = "YYYY-MM-DD"
		groupBy = "DATE(created_at)"
	case "yearly":
		dateFormat = "YYYY"
		groupBy = "TO_CHAR(created_at, 'YYYY')"
	default:
		dateFormat = "YYYY-MM"
		groupBy = "TO_CHAR(created_at, 'YYYY-MM')"
	}

	query := fmt.Sprintf(`SELECT TO_CHAR(created_at, '%s') as period,
		COUNT(*) as total, COALESCE(SUM(CASE WHEN status='completed' THEN total_amount ELSE 0 END),0),
		SUM(CASE WHEN status='completed' THEN 1 ELSE 0 END),
		SUM(CASE WHEN status='pending' THEN 1 ELSE 0 END),
		SUM(CASE WHEN status='cancelled' THEN 1 ELSE 0 END)
		FROM sales GROUP BY %s, TO_CHAR(created_at, '%s') ORDER BY period DESC LIMIT 12`, dateFormat, groupBy, dateFormat)

	rows, err := s.db.Query(query)
	if err != nil { return nil, err }
	defer rows.Close()

	var summaries []models.SalesSummary
	for rows.Next() {
		var sm models.SalesSummary
		rows.Scan(&sm.Period, &sm.TotalSales, &sm.TotalRevenue, &sm.CompletedSales, &sm.PendingSales, &sm.CancelledSales)
		if sm.CompletedSales > 0 { sm.AverageOrderValue = sm.TotalRevenue / float64(sm.CompletedSales) }
		summaries = append(summaries, sm)
	}
	return summaries, nil
}

func (s *ReportService) GetTopProducts(limit int) ([]models.TopProduct, error) {
	if limit < 1 || limit > 50 { limit = 10 }
	rows, err := s.db.Query(`SELECT si.product_id, COALESCE(p.name,'Unknown'), COALESCE(p.sku,''), COALESCE(p.category,''),
		SUM(si.quantity), SUM(si.line_total)
		FROM sale_items si LEFT JOIN products p ON si.product_id=p.id
		JOIN sales s ON si.sale_id=s.id WHERE s.status='completed'
		GROUP BY si.product_id, p.name, p.sku, p.category ORDER BY SUM(si.line_total) DESC LIMIT $1`, limit)
	if err != nil { return nil, err }
	defer rows.Close()

	var products []models.TopProduct
	for rows.Next() {
		var tp models.TopProduct
		rows.Scan(&tp.ProductID, &tp.ProductName, &tp.ProductSKU, &tp.Category, &tp.TotalQtySold, &tp.TotalRevenue)
		products = append(products, tp)
	}
	return products, nil
}

func (s *ReportService) GetRevenueReport() (*models.RevenueReport, error) {
	report := &models.RevenueReport{}
	s.db.QueryRow(`SELECT COALESCE(SUM(total_amount),0), COUNT(*) FROM sales WHERE status='completed'`).Scan(&report.TotalRevenue, &report.TotalTransactions)

	rows, _ := s.db.Query(`SELECT TO_CHAR(created_at,'YYYY-MM') as period, COALESCE(SUM(total_amount),0), COUNT(*)
		FROM sales WHERE status='completed' GROUP BY TO_CHAR(created_at,'YYYY-MM') ORDER BY period DESC LIMIT 12`)
	defer rows.Close()
	for rows.Next() {
		var pr models.PeriodRevenue
		rows.Scan(&pr.Period, &pr.Revenue, &pr.Count)
		report.RevenueByPeriod = append(report.RevenueByPeriod, pr)
	}
	return report, nil
}
