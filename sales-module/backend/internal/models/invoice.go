package models

import "time"

// Invoice status constants
const (
	InvoiceStatusDraft   = "draft"
	InvoiceStatusSent    = "sent"
	InvoiceStatusPaid    = "paid"
	InvoiceStatusOverdue = "overdue"
)

// Invoice represents a generated invoice
type Invoice struct {
	ID            int           `json:"id"`
	SaleID        *int          `json:"sale_id"`
	InvoiceNumber string        `json:"invoice_number"`
	CustomerID    *int          `json:"customer_id"`
	CustomerName  string        `json:"customer_name,omitempty"`
	IssueDate     string        `json:"issue_date"`
	DueDate       string        `json:"due_date"`
	Subtotal      float64       `json:"subtotal"`
	Discount      float64       `json:"discount"`
	Tax           float64       `json:"tax"`
	Total         float64       `json:"total"`
	Status        string        `json:"status"`
	Items         []InvoiceItem `json:"items,omitempty"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}

// InvoiceItem represents a line item on an invoice
type InvoiceItem struct {
	ID          int     `json:"id"`
	InvoiceID   int     `json:"invoice_id"`
	ProductID   *int    `json:"product_id"`
	Description string  `json:"description"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
	Discount    float64 `json:"discount"`
	LineTotal   float64 `json:"line_total"`
}

// InvoiceFilter holds filtering parameters for listing invoices
type InvoiceFilter struct {
	Status     string
	CustomerID int
	DateFrom   string
	DateTo     string
	Page       int
	Limit      int
}
