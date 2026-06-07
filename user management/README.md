# CRM System - User Management Module

## Project Description
This project implements the User Management Module for the Enterprise CRM System. It handles user authentication, registration, role management, and complete CRUD operations for managing system users. It is built using a modern alternative tech stack to achieve bonus marks:
- **Frontend**: Vue.js 3
- **Backend**: Go (Gin Framework)
- **Database**: PostgreSQL
- **Deployment**: Docker & Docker Compose

## Steps to Run
To run this project locally, ensure you have Docker and Docker Compose installed.

1. Clone or extract the repository.
2. Open a terminal in the project root directory.
3. Run the following command:
```bash
docker-compose up --build
```
4. The system will automatically build the images and start the database, backend, and frontend containers. No manual setup is required.

## Port Details
- **Frontend (Vue.js)**: `http://localhost` (Port 80)
- **Backend API (Go)**: `http://localhost:8080`
- **Database (PostgreSQL)**: `localhost:5432`

## API Endpoints (REST & JSON)

### Authentication
- `POST /api/v1/auth/register`: Register a new user
- `POST /api/v1/auth/login`: Authenticate user and receive JWT token
- `POST /api/v1/auth/logout`: Logout user

### User Management (Requires JWT Token)
- `GET /api/v1/users`: List all users
- `GET /api/v1/users/:id`: Get a specific user by ID
- `PUT /api/v1/users/:id`: Update user details (Role, Status, Email)
- `DELETE /api/v1/users/:id`: Delete a user

## Integration Details
This User Management module is designed to integrate with other modules in the CRM (like Sales, Customer Support) using shared identifiers.

- **Shared Identifier**: `user_id` (Primary Key). Other modules can store this `user_id` to track which sales representative handled a deal or which agent responded to a ticket.
- **REST APIs**: Other modules can call `GET /api/v1/users/:id` to retrieve details of a user handling a specific task.
- **Data Format**: All API responses use standard JSON format.
- **Authentication**: Other backend modules can verify the JWT tokens issued by this module's `/login` endpoint to securely identify the current logged-in user making cross-module requests.

## Database Initialization
The database tables and initial roles (Admin, Sales_Agent, Manager, Customer) are automatically created using the `init.sql` script mounted to the PostgreSQL container during startup.
