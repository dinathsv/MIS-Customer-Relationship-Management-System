package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

// ─── Models ──────────────────────────────────────────────────────────────────

type Feedback struct {
	FeedbackID string    `json:"feedback_id"`
	CustomerID string    `json:"customer_id"`
	Rating     int       `json:"rating"`
	Category   string    `json:"category"`
	Comments   string    `json:"comments"`
	CreatedAt  time.Time `json:"created_at"`
}

type FeedbackRequest struct {
	CustomerID string `json:"customer_id" binding:"required"`
	Rating     int    `json:"rating"      binding:"required,min=1,max=5"`
	Category   string `json:"category"    binding:"required"`
	Comments   string `json:"comments"`
}

type CategoryVolume struct {
	Category string  `json:"category"`
	Count    int     `json:"count"`
	AvgRating float64 `json:"avg_rating"`
}

type RatingDistribution struct {
	Rating int `json:"rating"`
	Count  int `json:"count"`
}

type AnalyticsSummary struct {
	TotalFeedback       int                  `json:"total_feedback"`
	AverageRating       float64              `json:"average_rating"`
	CategoryVolumes     []CategoryVolume     `json:"category_volumes"`
	RatingDistribution  []RatingDistribution `json:"rating_distribution"`
	RecentTrend         []DailyTrend         `json:"recent_trend"`
}

type DailyTrend struct {
	Date  string  `json:"date"`
	Count int     `json:"count"`
	Avg   float64 `json:"avg_rating"`
}

// ─── DB Connection ────────────────────────────────────────────────────────────

func connectDB() *sql.DB {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost port=5432 user=crmuser password=crmpassword dbname=crmdb sslmode=disable"
	}

	var db *sql.DB
	var err error

	// Retry logic for Docker startup ordering
	for i := 0; i < 10; i++ {
		db, err = sql.Open("postgres", dsn)
		if err == nil {
			if pingErr := db.Ping(); pingErr == nil {
				log.Println("✅ Database connected successfully")
				return db
			}
		}
		log.Printf("⏳ Waiting for database... attempt %d/10", i+1)
		time.Sleep(3 * time.Second)
	}
	log.Fatal("❌ Could not connect to database after 10 attempts:", err)
	return nil
}

// ─── Handlers ─────────────────────────────────────────────────────────────────

type Handler struct{ db *sql.DB }

func normalizeCustomerID(value string) string {
	raw := strings.TrimSpace(value)
	if raw == "" {
		return ""
	}
	return raw
}

// POST /api/feedback
func (h *Handler) SubmitFeedback(c *gin.Context) {
	var req FeedbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.CustomerID = normalizeCustomerID(req.CustomerID)
	if req.CustomerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "customer email is required"})
		return
	}

	validCategories := map[string]bool{
		"Product Quality": true, "Customer Support": true,
		"Delivery & Shipping": true, "Pricing & Value": true,
		"Website & App": true, "Returns & Refunds": true, "Other": true,
	}
	if !validCategories[req.Category] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category"})
		return
	}

	id := uuid.New().String()
	query := `
		INSERT INTO customer_feedback (feedback_id, customer_id, rating, category, comments)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING feedback_id, created_at`

	var fb Feedback
	err := h.db.QueryRow(query, id, req.CustomerID, req.Rating, req.Category, req.Comments).
		Scan(&fb.FeedbackID, &fb.CreatedAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save feedback"})
		log.Println("DB insert error:", err)
		return
	}

	fb.CustomerID = req.CustomerID
	fb.Rating = req.Rating
	fb.Category = req.Category
	fb.Comments = req.Comments

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Feedback submitted successfully",
		"feedback": fb,
	})
}

// GET /api/feedback
func (h *Handler) ListFeedback(c *gin.Context) {
	customerID := c.Query("customer_id")
	category := c.Query("category")
	limit := c.DefaultQuery("limit", "50")
	offset := c.DefaultQuery("offset", "0")

	query := `
		SELECT feedback_id, customer_id, rating, category, comments, created_at
		FROM customer_feedback
		WHERE ($1 = '' OR customer_id = $1)
		  AND ($2 = '' OR category = $2)
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4`

	rows, err := h.db.Query(query, customerID, category, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "query failed"})
		log.Println("DB query error:", err)
		return
	}
	defer rows.Close()

	feedbacks := []Feedback{}
	for rows.Next() {
		var fb Feedback
		var comments sql.NullString
		if err := rows.Scan(&fb.FeedbackID, &fb.CustomerID, &fb.Rating,
			&fb.Category, &comments, &fb.CreatedAt); err != nil {
			continue
		}
		if comments.Valid {
			fb.Comments = comments.String
		}
		feedbacks = append(feedbacks, fb)
	}

	// Total count for pagination
	var total int
	countQuery := `SELECT COUNT(*) FROM customer_feedback
		WHERE ($1 = '' OR customer_id = $1) AND ($2 = '' OR category = $2)`
	h.db.QueryRow(countQuery, customerID, category).Scan(&total)

	c.JSON(http.StatusOK, gin.H{
		"data":   feedbacks,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

// GET /api/analytics/summary
func (h *Handler) GetAnalyticsSummary(c *gin.Context) {
	summary := AnalyticsSummary{}

	// Total & average
	err := h.db.QueryRow(`SELECT COUNT(*), COALESCE(ROUND(AVG(rating)::numeric, 2), 0) FROM customer_feedback`).
		Scan(&summary.TotalFeedback, &summary.AverageRating)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "analytics query failed"})
		return
	}

	// Volume + avg rating by category
	catRows, err := h.db.Query(`
		SELECT category, COUNT(*) as count, ROUND(AVG(rating)::numeric, 2) as avg_rating
		FROM customer_feedback
		GROUP BY category
		ORDER BY count DESC`)
	if err == nil {
		defer catRows.Close()
		for catRows.Next() {
			var cv CategoryVolume
			catRows.Scan(&cv.Category, &cv.Count, &cv.AvgRating)
			summary.CategoryVolumes = append(summary.CategoryVolumes, cv)
		}
	}

	// Rating distribution (1-5)
	ratingRows, err := h.db.Query(`
		SELECT rating, COUNT(*) as count
		FROM customer_feedback
		GROUP BY rating
		ORDER BY rating ASC`)
	if err == nil {
		defer ratingRows.Close()
		// Initialise all 5 slots
		dist := map[int]int{1: 0, 2: 0, 3: 0, 4: 0, 5: 0}
		for ratingRows.Next() {
			var rd RatingDistribution
			ratingRows.Scan(&rd.Rating, &rd.Count)
			dist[rd.Rating] = rd.Count
		}
		for r := 1; r <= 5; r++ {
			summary.RatingDistribution = append(summary.RatingDistribution, RatingDistribution{Rating: r, Count: dist[r]})
		}
	}

	// Daily trend — last 14 days
	trendRows, err := h.db.Query(`
		SELECT
			TO_CHAR(created_at AT TIME ZONE 'UTC', 'YYYY-MM-DD') as day,
			COUNT(*) as count,
			ROUND(AVG(rating)::numeric, 2) as avg_rating
		FROM customer_feedback
		WHERE created_at >= NOW() - INTERVAL '14 days'
		GROUP BY day
		ORDER BY day ASC`)
	if err == nil {
		defer trendRows.Close()
		for trendRows.Next() {
			var dt DailyTrend
			trendRows.Scan(&dt.Date, &dt.Count, &dt.Avg)
			summary.RecentTrend = append(summary.RecentTrend, dt)
		}
	}

	if summary.CategoryVolumes == nil {
		summary.CategoryVolumes = []CategoryVolume{}
	}
	if summary.RatingDistribution == nil {
		summary.RatingDistribution = []RatingDistribution{}
	}
	if summary.RecentTrend == nil {
		summary.RecentTrend = []DailyTrend{}
	}

	c.JSON(http.StatusOK, summary)
}

// GET /api/feedback/:id — fetch single feedback by CRM customer_id
func (h *Handler) GetFeedbackByCustomer(c *gin.Context) {
	customerID := c.Param("customer_id")
	rows, err := h.db.Query(`
		SELECT feedback_id, customer_id, rating, category, comments, created_at
		FROM customer_feedback WHERE customer_id = $1
		ORDER BY created_at DESC`, customerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "query failed"})
		return
	}
	defer rows.Close()

	feedbacks := []Feedback{}
	for rows.Next() {
		var fb Feedback
		var comments sql.NullString
		rows.Scan(&fb.FeedbackID, &fb.CustomerID, &fb.Rating, &fb.Category, &comments, &fb.CreatedAt)
		if comments.Valid {
			fb.Comments = comments.String
		}
		feedbacks = append(feedbacks, fb)
	}
	c.JSON(http.StatusOK, gin.H{"customer_id": customerID, "feedback": feedbacks, "count": len(feedbacks)})
}

// ─── Health check ─────────────────────────────────────────────────────────────

func (h *Handler) Health(c *gin.Context) {
	if err := h.db.Ping(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unhealthy", "db": "unreachable"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "healthy", "service": "feedback-api", "version": "1.0.0"})
}

// ─── Main ─────────────────────────────────────────────────────────────────────

func main() {
	db := connectDB()
	defer db.Close()

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	h := &Handler{db: db}

	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	// CORS — allow Vue dev server and production frontend
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	// Routes
	r.GET("/health", h.Health)

	api := r.Group("/api")
	{
		api.POST("/feedback", h.SubmitFeedback)
		api.GET("/feedback", h.ListFeedback)
		api.GET("/feedback/customer/:customer_id", h.GetFeedbackByCustomer)
		api.GET("/analytics/summary", h.GetAnalyticsSummary)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Feedback API running on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
