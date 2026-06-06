package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

type Complaint struct {
	ID          int    `json:"id"`
	CustomerID  int    `json:"customer_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
}

var db *sql.DB

func main() {
	var err error
	// Fetch connection string from environment variables
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://user:password@localhost:5432/crm_db?sslmode=disable"
	}

	db, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// API Routing
	http.HandleFunc("/api/complaints", handleComplaints)

	fmt.Println("Backend server running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleComplaints(w http.ResponseWriter, r *http.Request) {
	// Simple CORS handling for integration and local testing
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if r.Method == "GET" {
		rows, err := db.Query("SELECT id, customer_id, title, description, status, created_at FROM complaints")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var complaints []Complaint
		for rows.Next() {
			var c Complaint
			err := rows.Scan(&c.ID, &c.CustomerID, &c.Title, &c.Description, &c.Status, &c.CreatedAt)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			complaints = append(complaints, c)
		}

		json.NewEncoder(w).Encode(complaints)

	} else if r.Method == "POST" {
		var c Complaint
		err := json.NewDecoder(r.Body).Decode(&c)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = db.QueryRow(
			"INSERT INTO complaints (customer_id, title, description, status) VALUES ($1, $2, $3, 'Pending') RETURNING id",
			c.CustomerID, c.Title, c.Description,
		).Scan(&c.ID)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(c)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
