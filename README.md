# MIS Enterprise CRM Platform

A comprehensive, microservice-based Enterprise CRM system integrating four distinct business modules into a unified Docker-based architecture.

## 📌 Project Description

This system is designed to handle end-to-end CRM operations. It consists of four independent modules that communicate over a shared internal Docker network (`crm-network`), while a centralized Enterprise Portal provides unified access to all interfaces.

### The 4 Modules:
1. **Sales Module**: Manages product inventory, generates invoices, and handles sales tracking.
2. **User Management Module**: Administers staff accounts, handles JWT-based authentication, and manages role-based access control.
3. **Complain Management Module**: Tracks customer issues, assigns statuses, and monitors resolution progress.
4. **Customer Feedback Module**: Collects customer ratings, categorizes feedback, and provides analytics dashboards.

All modules are completely dockerized with their own dedicated PostgreSQL databases, Golang backend APIs, and Vue.js/Vite frontend SPAs served via Nginx.

---

## 🚀 Steps to Run the System

This project is built for zero manual setup. It runs entirely inside Docker.

### Prerequisites
- Docker and Docker Compose must be installed on your machine.

### Deployment Instructions
1. Open your terminal and navigate to the root directory of this project (where the `docker-compose.yml` file is located).
2. Execute the following command to build and spin up all 13 containers:
   ```bash
   docker-compose up --build
   ```
3. Wait for the build process to complete and the databases to initialize.
4. Open your web browser and navigate to the Central Enterprise Portal at: **http://localhost:3000**

To stop the system gracefully, run: `docker-compose down`

---

## 🌐 Port Details & Access Links

| Module | Frontend UI | Backend API | Database (PostgreSQL) | Default Admin Credentials |
|--------|------------|-------------|----------------------|---------------------------|
| **Central Portal** | `http://localhost:3000` | - | - | - |
| **Sales Module** | `http://localhost:8086` | `http://localhost:8080` | `localhost:5436` | `admin` / `admin123` |
| **Complain Mgmt** | `http://localhost:3001` | `http://localhost:8081` | `localhost:5433` | `admin` / `admin123` |
| **User Mgmt** | `http://localhost:8083` | `http://localhost:8082` | `localhost:5434` | `admin@gmail.com` / `admin123` |
| **Feedback Module**| `http://localhost:8085` | `http://localhost:8084` | `localhost:5435` | _(Public Submission Form)_ |

---

## 🔗 API Endpoints

The system exposes RESTful APIs for each module. A complete Postman Collection is included in the project root (`MIS_Enterprise_API_Collection.json`) for easy testing. Below is a high-level summary:

### User Management API (`:8082`)
- `POST /api/v1/auth/login` - Authenticate user
- `POST /api/v1/auth/register` - Register new user
- `GET /api/v1/users` - Get all users (Admin only)
- `PUT /api/v1/users/me/profile` - Update own profile

### Sales API (`:8080`)
- `GET /api/v1/sales` - List all sales
- `POST /api/v1/sales` - Create a new sale
- `GET /api/v1/products` - List available products
- `POST /api/v1/invoices/generate/:saleId` - Generate an invoice

### Complain Management API (`:8081`)
- `GET /api/complaints` - Fetch all complaints
- `POST /api/complaints` - Submit a new complaint

### Feedback API (`:8084`)
- `POST /api/feedback` - Submit customer feedback
- `GET /api/feedback` - List all feedback
- `GET /api/analytics/summary` - View feedback analytics

---

## ⚙️ Integration Details

1. **Shared Network (`crm-network`)**: All containers are bridged together on a custom Docker network. This allows backends to communicate securely with their respective databases via internal DNS (e.g., the User Management backend connects to `usermgmt_db:5432` rather than localhost).
2. **Global Navigation**: A universal floating "Module Switcher" is injected into the DOM of every frontend. This allows users to seamlessly jump between the Sales, Complain, User, and Feedback modules without returning to the main portal.
3. **Automated Seed Data**: All 4 PostgreSQL databases have dedicated `init.sql` scripts mounted to `/docker-entrypoint-initdb.d/`. Upon the first `docker-compose up`, these scripts automatically create tables, insert default admin users, and populate sample analytics data.

---

## 🖥️ Demo Instructions

As per the demo requirements, the system hosts the full ecosystem on a single localhost environment.

1. Ensure no other services (like local Postgres) are running on ports `80, 3000, 8080-8086, 5433-5436`.
2. Run `docker-compose up --build`.
3. Start the demo by opening `http://localhost:3000`. This showcases the **Enterprise Portal** bridging all group modules.
4. Click through to each module to demonstrate its frontend capabilities.
5. Use the floating "Switch Module" button in the bottom right corner to smoothly transition between groups during the presentation.
