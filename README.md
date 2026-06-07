# 🌌 MIS Enterprise CRM Platform

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://golang.org)
[![Vue Version](https://img.shields.io/badge/Vue.js-3.x-4FC08D?style=for-the-badge&logo=vue.js&logoColor=white)](https://vuejs.org)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15%2B-4169E1?style=for-the-badge&logo=postgresql&logoColor=white)](https://www.postgresql.org)
[![Docker](https://img.shields.io/badge/Docker-Enabled-2496ED?style=for-the-badge&logo=docker&logoColor=white)](https://www.docker.com)
[![Architecture](https://img.shields.io/badge/Architecture-Microservices-orange?style=for-the-badge)](https://microservices.io)

A highly responsive, production-ready, microservices-driven Enterprise CRM system. Featuring containerized independent modules communicating via a secure internal bridge network, standardizing user experience with a unified Glassmorphic UI design.

---

## 🗺️ System Architecture

The ecosystem consists of **4 core business modules** and **1 central gateway/portal**. All services run isolated within Docker, leveraging standard seed databases initialized automatically on first spin-up.

```mermaid
graph TD
    CP[Central Portal - :3000] --> Sales[Sales Module - :8086]
    CP --> Complain[Complain Module - :3001]
    CP --> UserMgmt[User Management - :8083]
    CP --> Feedback[Feedback Module - :8085]
    
    subgraph Microservices Network (crm-network)
        Sales -.-> SalesAPI[Sales API - :8080] --> SalesDB[(Sales DB - :5436)]
        Complain -.-> ComplainAPI[Complain API - :8081] --> ComplainDB[(Complain DB - :5433)]
        UserMgmt -.-> UserAPI[User API - :8082] --> UserDB[(User DB - :5434)]
        Feedback -.-> FeedbackAPI[Feedback API - :8084] --> FeedbackDB[(Feedback DB - :5435)]
    end
```

---

## ✨ Key Features & UX Standards

- **Unified Glassmorphic Header**: Every module features an elegant, responsive top navigation bar complete with custom branding, active module routing status, real-time user identity (Role & Username), and direct Logout controls.
- **Enterprise Portal Switcher**: Seamlessly switch between micro-frontends with a single click, either through the Top Portal button or the floating Module Switcher.
- **Role-Based Access Control (RBAC)**:
  - **Standard Users**: Can update their personal profiles and view features dedicated to their respective departments.
  - **Administrators**: Get advanced dashboard management views showing all registered users system-wide.
- **Microservices Isolation**: Fully containerized Go backends, Vite-powered Vue frontends, and distinct PostgreSQL databases initialized dynamically via `init.sql` seed scripts.

---

## 🚀 Getting Started (Run with Zero Setup)

This platform is configured for instant boot. Make sure Docker is running on your machine before starting.

### Prerequisites
- [Docker Desktop](https://www.docker.com/products/docker-desktop/) or Docker Engine
- [Docker Compose](https://docs.docker.com/compose/)

### Execution Command
Navigate to the root directory containing the `docker-compose.yml` file and run:

```bash
docker-compose up --build
```

> [!NOTE]
> *On the first run, the system automatically pulls base images, builds frontends, compiles Go backends, creates databases, and runs SQL seeds (located in individual db folders). This can take a few minutes.*

To tear down the containers and release resources, run:
```bash
docker-compose down
```

---

## 🌐 Modules & Port Registry

Access the modules directly using the following configurations:

| Module / Component | Frontend Port | Backend API | Database | Default Credentials |
| :--- | :---: | :---: | :---: | :--- |
| **🌌 Central Portal** | [`:3000`](http://localhost:3000) | — | — | *Landing directory for all modules* |
| **💼 Sales Module** | [`:8086`](http://localhost:8086) | `:8080` | `:5436` | `admin` / `admin123` |
| **🛡️ Complain Mgmt** | [`:3001`](http://localhost:3001) | `:8081` | `:5433` | `admin` / `admin123` |
| **👥 User Management** | [`:8083`](http://localhost:8083) | `:8082` | `:5434` | `admin@gmail.com` / `admin123` |
| **📝 Feedback Module** | [`:8085`](http://localhost:8085) | `:8084` | `:5435` | *Public Feedback Submission Form* |

---

## 🔌 API Endpoints Summary

A pre-configured Postman Collection is available in the root folder: **`MIS_Enterprise_API_Collection.json`**. Below is a quick overview:

### 👥 User Management API (`:8082`)
* `POST /api/v1/auth/login` — Sign in and obtain JWT
* `POST /api/v1/auth/register` — Standard registration
* `GET /api/v1/users` — Fetch all system users *(Admin Only)*
* `PUT /api/v1/users/me/profile` — Update active user profile details

### 💼 Sales API (`:8080`)
* `GET /api/v1/sales` — Fetch entire sales logs
* `POST /api/v1/sales` — Insert a new sales record
* `GET /api/v1/products` — Retrieve inventory listing
* `POST /api/v1/invoices/generate/:saleId` — Create formatted invoice PDF/metadata

### 🛡️ Complain API (`:8081`)
* `GET /api/complaints` — Retrieve complaints list
* `POST /api/complaints` — Submit a customer service complaint

### 📝 Feedback API (`:8084`)
* `POST /api/feedback` — Submit feedback (linked with Customer ID)
* `GET /api/feedback` — Fetch customer feedback entries
* `GET /api/analytics/summary` — Fetch computed analytics & ratings count

---

## ⚙️ Core Technical Integrations

1. **Shared Docker Network (`crm-network`)**: Custom internal bridge network. Backends query database service nodes via internal DNS names (e.g. `usermgmt_db:5432`) rather than localhost, preventing port conflicts inside Docker.
2. **PostgreSQL Automounts**: Databases utilize `init.sql` configurations mounted at `/docker-entrypoint-initdb.d/`. Tables, constraints, indices, and default user accounts are fully seeded immediately on startup.
3. **Micro-frontend Styling**: Implements premium Glassmorphism design system using vanilla CSS rules, custom variables (`var(--glass-bg)` etc.), dynamic backdrops, and modern typography.
