# GEMINI.md - Delivery App Backend

## Project Overview

This project is the backend for a full-featured delivery platform built with Golang and the Gin framework. It supports real-time order tracking, chat functionality between customers and shippers, and an admin dashboard for management. The architecture is layered, separating concerns into handlers, services, and repositories.

**Key Technologies:**

*   **Backend:** Golang (Gin Framework)
*   **Database:** MySQL
*   **Authentication:** JWT with Email OTP verification
*   **Real-time Features:** WebSocket for chat and location updates
*   **File Storage:** Cloudinary for image uploads
*   **Email:** SMTP for OTP and notifications

## Building and Running

### Prerequisites

*   Go (version 1.22 or higher recommended)
*   MySQL
*   An SMTP server for email
*   A Cloudinary account for file storage

### Local Development

1.  **Clone the repository.**

2.  **Set up environment variables:**
    Create a `.env` file in the root directory and add the following, replacing the placeholder values:

    ```bash
    # Database
    DB_URL="user:password@tcp(127.0.0.1:3306)/DeliveryAppDB?charset=utf8mb4&parseTime=True&loc=Local"

    # JWT
    JWT_SECRET="your_jwt_secret"

    # Email (SMTP)
    EMAIL_HOST="smtp.example.com"
    EMAIL_PORT=587
    EMAIL_USER="your_email@example.com"
    EMAIL_PASSWORD="your_email_password"
    EMAIL_SENDER="Your App Name <no-reply@example.com>"

    # Cloudinary
    CLOUDINARY_URL="cloudinary://api_key:api_secret@cloud_name"
    ```

3.  **Install dependencies:**

    ```bash
    go mod tidy
    ```

4.  **Run the database schema:**
    The database schema is located at `database/schema.sql`. You can import this file into your MySQL database to create the necessary tables.

5.  **Run the application:**

    ```bash
    go run cmd/main.go
    ```

    The server will start on `http://localhost:8080`.

### Testing

There are no specific testing instructions in the project. However, you can add tests in the `test` directory and run them using the standard Go testing tools.

```bash
go test ./...
```

## Development Conventions

*   **Layered Architecture:** The project follows a clean, layered architecture:
    *   `handlers`: Responsible for handling HTTP requests and responses.
    *   `services`: Contains the business logic of the application.
    *   `repository`: Handles database operations.
    *   `dto`: Data Transfer Objects for request and response validation.
    *   `models`: Represents the database entities.
*   **Dependency Injection:** Dependencies are injected from the `main` function, promoting loose coupling and testability.
*   **Routing:** All API routes are defined in the `routes/routes.go` file.
*   **Middleware:** Middleware is used for authentication (`AuthMiddleware`) and role-based access control (`RoleMiddleWare`).
*   **Configuration:** All configuration is loaded from environment variables using the `godotenv` package.
*   **Database:** The project uses a `Unit of Work` pattern for database transactions.

## API Endpoints

A summary of the main API endpoints can be found in `routes/routes.go`. The base path for all API routes is `/api/v1`.

### Authentication

*   `POST /signup`
*   `POST /verify-otp`
*   `POST /login`
*   `POST /forgot-password`
*   `POST /reset-password`
*   `POST /refresh-access-token`
*   `POST /logout`

### Protected Routes

These routes require a valid JWT token in the `Authorization` header.

*   **Customer:** `/customer/*`
*   **Shipper:** `/shipper/*`
*   **Admin:** `/admin/*`

### WebSocket

*   `GET /ws`: Establishes a WebSocket connection for real-time communication.
