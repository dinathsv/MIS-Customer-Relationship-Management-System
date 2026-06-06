# 📊 Customer Feedback Analytical Module
### CRM Enterprise System — Continuous Assessment

A production-ready, fully containerised **Customer Feedback Analytics** module built with **Vue 3 + Go + PostgreSQL**, featuring a glassmorphism UI dashboard and RESTful API.

---

## 🗂️ Project Structure

```
feedback-module/
├── docker-compose.yml          # Orchestrates all 3 services
├── database/
│   └── init.sql                # Schema + 40 mock data rows
├── backend/
│   ├── Dockerfile              # Multi-stage Go build
│   ├── go.mod / go.sum
│   └── main.go                 # Gin REST API + analytics engine
└── frontend/
    ├── Dockerfile              # Node build → Nginx serve
    ├── nginx.conf              # SPA + reverse proxy config
    ├── package.json
    ├── vite.config.js
    ├── tailwind.config.js
    └── src/
        ├── App.vue             # Root layout, navigation
        ├── main.js
        ├── assets/main.css     # Glassmorphism + Tailwind
        ├── composables/
        │   └── useApi.js       # Axios API client
        └── components/
            ├── FeedbackForm.vue        # Submission form
            └── AnalyticsDashboard.vue  # Charts & KPIs
```

---

## 🚀 Quick Start — Docker (Zero Setup)

**Prerequisites:** Docker Desktop (or Docker + Docker Compose v2)

```bash
# 1. Clone / unzip the project
cd feedback-module

# 2. Build and launch all services
docker-compose up --build

# 3. Open the app
open http://localhost:3000
```

That's it. The database is seeded automatically with 40 mock rows on first boot.

> **Tip:** On first run the backend waits for PostgreSQL to be healthy before starting (retry logic built in). Allow ~30 seconds for everything to be ready.

---

## 🌐 Port Reference

| Service    | Container Port | Host Port | URL                        |
|------------|---------------|-----------|----------------------------|
| Frontend   | 80            | **3000**  | http://localhost:3000      |
| Backend    | 8080          | **8080**  | http://localhost:8080      |
| PostgreSQL | 5432          | **5432**  | localhost:5432             |

---

## 🗄️ Database Schema

```sql
CREATE TABLE customer_feedback (
    feedback_id   UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    customer_id   VARCHAR(64)  NOT NULL,   -- CRM integration key
    rating        SMALLINT     NOT NULL CHECK (rating BETWEEN 1 AND 5),
    category      VARCHAR(64)  NOT NULL,
    comments      TEXT,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
```

**Valid categories:** `Product Quality`, `Customer Support`, `Delivery & Shipping`, `Pricing & Value`, `Website & App`, `Returns & Refunds`, `Other`

**Indexes:** `customer_id`, `category`, `rating`, `created_at` (optimised for CRM queries)

---

## 📡 API Documentation

Base URL: `http://localhost:8080`

---

### `GET /health`
Service health check.

**Response `200`:**
```json
{
  "status": "healthy",
  "service": "feedback-api",
  "version": "1.0.0"
}
```

---

### `POST /api/feedback`
Submit a new customer feedback entry. The `customer_id` field links this record to the CRM customer profile.

**Request Body:**
```json
{
  "customer_id": "CRM-2001",
  "rating": 5,
  "category": "Product Quality",
  "comments": "Outstanding build quality. Exceeded all expectations."
}
```

| Field         | Type    | Required | Validation                    |
|---------------|---------|----------|-------------------------------|
| `customer_id` | string  | ✅       | Non-empty, links to CRM       |
| `rating`      | integer | ✅       | 1–5 inclusive                 |
| `category`    | string  | ✅       | Must be a valid category      |
| `comments`    | string  | ❌       | Optional free-text            |

**Response `201`:**
```json
{
  "message": "Feedback submitted successfully",
  "feedback": {
    "feedback_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
    "customer_id": "CRM-2001",
    "rating": 5,
    "category": "Product Quality",
    "comments": "Outstanding build quality.",
    "created_at": "2025-06-15T10:30:00Z"
  }
}
```

**Response `400` (validation error):**
```json
{ "error": "invalid category" }
```

---

### `GET /api/feedback`
List all feedback with optional filtering and pagination.

**Query Parameters:**

| Parameter     | Type    | Description                        |
|---------------|---------|------------------------------------|
| `customer_id` | string  | Filter by CRM customer ID          |
| `category`    | string  | Filter by category name            |
| `limit`       | integer | Max records to return (default 50) |
| `offset`      | integer | Pagination offset (default 0)      |

**Examples:**
```bash
# All feedback
GET /api/feedback

# Filter by customer (CRM integration)
GET /api/feedback?customer_id=CRM-1001

# Filter by category with pagination
GET /api/feedback?category=Customer+Support&limit=10&offset=0
```

**Response `200`:**
```json
{
  "data": [
    {
      "feedback_id": "3fa85f64-...",
      "customer_id": "CRM-1001",
      "rating": 5,
      "category": "Product Quality",
      "comments": "Absolutely love the build quality.",
      "created_at": "2025-06-14T09:00:00Z"
    }
  ],
  "total": 40,
  "limit": "50",
  "offset": "0"
}
```

---

### `GET /api/feedback/customer/:customer_id`
Retrieve all feedback for a specific CRM customer. This is the primary CRM integration endpoint — look up a customer in the CRM and pass their ID here.

**Example:**
```bash
GET /api/feedback/customer/CRM-1001
```

**Response `200`:**
```json
{
  "customer_id": "CRM-1001",
  "count": 1,
  "feedback": [
    {
      "feedback_id": "3fa85f64-...",
      "customer_id": "CRM-1001",
      "rating": 5,
      "category": "Product Quality",
      "comments": "Absolutely love the build quality.",
      "created_at": "2025-06-14T09:00:00Z"
    }
  ]
}
```

---

### `GET /api/analytics/summary`
Returns aggregated analytics for the CRM dashboard: average rating, volume by category, rating distribution, and 14-day trend.

**Example:**
```bash
GET /api/analytics/summary
```

**Response `200`:**
```json
{
  "total_feedback": 40,
  "average_rating": 3.55,
  "category_volumes": [
    { "category": "Product Quality",  "count": 6, "avg_rating": 3.83 },
    { "category": "Customer Support", "count": 6, "avg_rating": 3.33 }
  ],
  "rating_distribution": [
    { "rating": 1, "count": 4 },
    { "rating": 2, "count": 6 },
    { "rating": 3, "count": 7 },
    { "rating": 4, "count": 11 },
    { "rating": 5, "count": 12 }
  ],
  "recent_trend": [
    { "date": "2025-06-14", "count": 14, "avg_rating": 3.71 },
    { "date": "2025-06-13", "count": 10, "avg_rating": 3.40 }
  ]
}
```

---

## 🎨 Frontend Features

- **Glassmorphism UI** — frosted glass cards, `backdrop-filter: blur`, semi-transparent backgrounds
- **Animated gradient background** — multi-layer radial gradients with CSS `hue-rotate` animation
- **Feedback Form** — star rating picker, category select, CRM customer ID input, live validation
- **Analytics Dashboard** — 4 KPI stat cards, bar chart (category volume), doughnut chart (rating distribution), line chart (14-day trend), category performance table with progress bars
- **Auto-refresh** — after form submission, auto-switches to analytics and reloads data
- **API status indicator** — live health check shown in the header
- **Responsive** — works on mobile, tablet, and desktop

---

## 🛠️ Tech Stack

| Layer      | Technology                                  |
|------------|---------------------------------------------|
| Frontend   | Vue 3 (Composition API), Vite, Tailwind CSS |
| Charts     | Chart.js + vue-chartjs (Bar, Doughnut, Line)|
| Backend    | Go 1.22, Gin framework                      |
| Database   | PostgreSQL 16                               |
| Deployment | Docker, Docker Compose, Nginx               |

---

## 🔧 Development (Without Docker)

**Backend:**
```bash
cd backend
export DATABASE_URL="host=localhost port=5432 user=crmuser password=crmpassword dbname=crmdb sslmode=disable"
go run main.go
```

**Frontend:**
```bash
cd frontend
npm install
npm run dev   # → http://localhost:5173
```

**Database:**
```bash
psql -U crmuser -d crmdb -f database/init.sql
```

---

## 🐳 Useful Docker Commands

```bash
# Start in background
docker-compose up -d --build

# View logs
docker-compose logs -f backend
docker-compose logs -f frontend
docker-compose logs -f db

# Stop all services
docker-compose down

# Full reset (removes DB volume)
docker-compose down -v

# Rebuild single service
docker-compose up -d --build backend
```

---

## 🔑 CRM Integration Notes

The `customer_id` field is the bridge between this feedback module and the broader CRM system:

1. When a CRM user views a customer profile, the CRM queries `GET /api/feedback/customer/{crm_customer_id}`
2. The feedback module returns all historical feedback for that customer
3. New feedback submitted via the form is immediately queryable by the same endpoint
4. The analytics summary (`GET /api/analytics/summary`) aggregates across all customers for manager dashboards

---

*Built for CRM Enterprise Systems — Continuous Assessment*
