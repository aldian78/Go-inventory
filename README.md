# Smart Inventory Core System

This is a backend microservice for a Smart Inventory System built with Go (Golang) and Fiber, using PostgreSQL as the database.

## Prerequisites

- Go 1.20+
- Docker & Docker Compose (optional, for running PostgreSQL)
- PostgreSQL (if not using Docker)

## Setup

1.  **Clone the repository**
2.  **Database Setup**
    *   You can use the provided `docker-compose.yml` to start a PostgreSQL instance:
        ```bash
        docker-compose up -d
        ```
    *   Or configure your local PostgreSQL and update the `.env` file in the `backend` directory.

3.  **Backend Setup**
    *   Navigate to the `backend` directory:
        ```bash
        cd backend
        ```
    *   Install dependencies:
        ```bash
        go mod tidy
        ```
    *   Run the application:
        ```bash
        go run main.go
        ```
    *   The server will start on `http://localhost:8080`.

## API Endpoints

### Products
*   `GET /api/products`: List all products.
*   `POST /api/products`: Create a new product.
*   `PUT /api/products/:id/adjust`: Adjust stock (Physical & Available).

### Stock In
*   `POST /api/stock-in`: Create a Stock In request (Status: CREATED).
*   `PUT /api/stock-in/:id/status`: Update status (CREATED -> IN_PROGRESS -> DONE).
    *   **Note**: Physical stock increases only when status becomes DONE.
*   `GET /api/stock-in`: Get Stock In logs.

### Stock Out (Two-Phase Commitment)
*   `POST /api/stock-out/allocate`: Stage 1 - Allocation (Status: ALLOCATED). Checks available stock and reserves it.
*   `PUT /api/stock-out/:id/execute`: Stage 2 - Execution (Status: IN_PROGRESS).
*   `PUT /api/stock-out/:id/complete`: Complete the process (Status: DONE). Deducts physical stock.
*   `PUT /api/stock-out/:id/cancel`: Cancel the request. Rollbacks allocated stock if not DONE.
*   `GET /api/stock-out`: Get Stock Out logs.

### Reports
*   `GET /api/reports`: Get transaction reports for completed (DONE) Stock In and Stock Out transactions.

## Project Structure

*   `backend/`: Go backend code.
    *   `config/`: Database configuration.
    *   `controllers/`: Request handlers.
    *   `models/`: Database models.
    *   `routes/`: API route definitions.
*   `docker-compose.yml`: Docker configuration for PostgreSQL.

## Features Implemented

*   **Stock In**: Tracks status changes. Updates physical stock only on DONE.
*   **Inventory**: Tracks Physical Stock and Available Stock separately.
*   **Stock Out**: Implements Two-Phase Commitment (Allocation -> Execution -> Done). Handles cancellations and rollbacks.
*   **Reports**: Generates reports for completed transactions.
