package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
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

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"` // "admin" or "user"
}

var db *sql.DB
var JWTKey = []byte("CRM_SUPER_SECRET_KEY_2026")

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

// CORS Handler Wrapper
func corsHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

// Get token from Authorization header
func getTokenFromRequest(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}
	parts := strings.Split(authHeader, " ")
	if len(parts) == 2 && parts[0] == "Bearer" {
		return parts[1]
	}
	return ""
}

func getUserFromToken(tokenString string) *User {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return JWTKey, nil
	})

	if err != nil || !token.Valid {
		return nil
	}

	roleIDFloat, ok := claims["role_id"].(float64)
	if !ok {
		return nil
	}

	roleID := int(roleIDFloat)
	username, _ := claims["username"].(string)
	userIDFloat, _ := claims["user_id"].(float64)

	role := "user"
	if roleID == 1 {
		role = "admin"
	}

	return &User{
		ID:       int(userIDFloat),
		Username: username,
		Role:     role,
	}
}

func handleComplaints(w http.ResponseWriter, r *http.Request) {
	// CORS handling
	corsHeaders(w)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// Check for ID in query string for PUT/DELETE operations
	id := r.URL.Query().Get("id")

	if r.Method == "GET" {
		rows, err := db.Query("SELECT id, customer_id, title, description, status, created_at FROM complaints")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		complaints := []Complaint{}
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

	} else if r.Method == "PUT" {
		if id == "" {
			http.Error(w, "ID is required", http.StatusBadRequest)
			return
		}

		// Check authentication and admin role
		token := getTokenFromRequest(r)
		user := getUserFromToken(token)
		if user == nil || user.Role != "admin" {
			http.Error(w, "Unauthorized - Admin access required", http.StatusUnauthorized)
			return
		}

		var updateData map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&updateData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Update the complaint with new status
		status, ok := updateData["status"].(string)
		if !ok {
			http.Error(w, "Status field is required", http.StatusBadRequest)
			return
		}

		_, err = db.Exec("UPDATE complaints SET status = $1 WHERE id = $2", status, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Fetch and return the updated complaint
		var c Complaint
		err = db.QueryRow("SELECT id, customer_id, title, description, status, created_at FROM complaints WHERE id = $1", id).Scan(
			&c.ID, &c.CustomerID, &c.Title, &c.Description, &c.Status, &c.CreatedAt,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(c)

	} else if r.Method == "DELETE" {
		if id == "" {
			http.Error(w, "ID is required", http.StatusBadRequest)
			return
		}

		// Check authentication and admin role
		token := getTokenFromRequest(r)
		user := getUserFromToken(token)
		if user == nil || user.Role != "admin" {
			http.Error(w, "Unauthorized - Admin access required", http.StatusUnauthorized)
			return
		}

		result, err := db.Exec("DELETE FROM complaints WHERE id = $1", id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if rowsAffected == 0 {
			http.Error(w, "Complaint not found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Complaint deleted successfully"})

	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
