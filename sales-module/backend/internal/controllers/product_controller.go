package controllers

import (
	"net/http"
	"sales-module/internal/models"
	"sales-module/internal/services"

	"github.com/gin-gonic/gin"
)

type ProductController struct {
	service *services.ProductService
}

func NewProductController(service *services.ProductService) *ProductController {
	return &ProductController{service: service}
}

func (pc *ProductController) List(c *gin.Context) {
	products, err := pc.service.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, products)
}

func (pc *ProductController) Create(c *gin.Context) {
	var input struct {
		Name        string  `json:"name" binding:"required"`
		SKU         string  `json:"sku" binding:"required"`
		Category    string  `json:"category"`
		Description string  `json:"description"`
		Price       float64 `json:"price" binding:"required"`
		StockQty    int     `json:"stock_qty"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product := &models.Product{
		Name:        input.Name,
		SKU:         input.SKU,
		Category:    input.Category,
		Description: input.Description,
		Price:       input.Price,
		StockQty:    input.StockQty,
	}

	if err := pc.service.Create(product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, product)
}
