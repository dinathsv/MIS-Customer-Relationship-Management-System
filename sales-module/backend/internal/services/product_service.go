package services

import (
	"database/sql"
	"sales-module/internal/models"
)

type ProductService struct {
	db *sql.DB
}

func NewProductService(db *sql.DB) *ProductService {
	return &ProductService{db: db}
}

func (s *ProductService) List() ([]models.Product, error) {
	rows, err := s.db.Query("SELECT id, name, sku, category, price, stock_qty FROM products ORDER BY name ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.SKU, &p.Category, &p.Price, &p.StockQty); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func (s *ProductService) Create(p *models.Product) error {
	query := `INSERT INTO products (name, sku, category, description, price, stock_qty) 
	          VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	return s.db.QueryRow(query, p.Name, p.SKU, p.Category, p.Description, p.Price, p.StockQty).Scan(&p.ID)
}
