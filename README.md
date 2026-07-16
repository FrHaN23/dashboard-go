# Dashboard Go - Article REST API

[![Go Version](https://img.shields.io/badge/Go-1.22%2B-00ADD8?style=flat&logo=go)](https://img.shields.io/badge/Docker-Enabled-2496ED?style=flat&logo=docker)](https://www.docker.com/)

A robust, production-ready REST API for managing articles (*Article Dashboard*), built purely using the **Go Standard Library** (without heavy third-party web frameworks like Fiber or Gin) and backed by a **MySQL 8.0** database. This project is structured following the project structured following the principles of **Clean Architecture** and features dynamic pagination, sorting, soft deletes, and CORS handling.

## 🚀 Key Features

- **Clean Architecture & SOLID Principles**: Clear separation of layers between *Driver/Server*, *Handler*, *Usecase*, and *Repository*.
- **Docker Compose Orchestration**: Instantly spins up the MySQL database and the Go application (using an optimized *multi-stage build* Dockerfile).
- **Go 1.22+ Native Routing**: Leverages the newly improved `http.ServeMux` for native HTTP method routing.
- **Graceful Shutdown**: Safely shuts down the server without interrupting active network connections (*zero-downtime shutdown*).
- **Soft Delete (Trash)**: Deletes articles by moving them to a `thrash` status instead of permanently wiping them from the database.
- **Pagination & Sorting**: Returns comprehensive pagination metadata (`total_pages`, `total_rows`) and automatically sorts data by the latest `updated_date` for smooth frontend integration.
- **CORS Middleware**: Native middleware built to handle CORS safely (defaults to `http://localhost:3000` for frontend development).

---

## 🛠️ Prerequisites

Ensure you have the following installed on your machine:
- [Docker & Docker Compose](https://www.docker.com/products/docker-desktop/)
- [Go 1.22+](https://go.dev/dl/) (Optional: only if you wish to run the app natively outside Docker)

---

## 🏃 How to Run the Application

You can run the application using **Docker Compose** (recommended for a quick setup) or **Locally** on your machine.

### Method A: Using Docker Compose (Quick & Easy)

This method automatically builds the Go app, starts MySQL, configures user privileges, and hooks the method automatically builds the Go app, starts MySQL, configures user privileges, and hooks them together inside an isolated bridge network.

1. **Clone the repository:**
   ```bash
   git clone https://github.com/FrHaN23/dashboard-go.git
   cd dashboard-go
   ```

2. **Start the Docker services:**
   ```bash
   docker compose up --build -d
   ```
   *This starts the containers in detached mode.*

3. **Initialize & Seed the Database:**
   Run the following command in your terminal to create the schema and populate the database with seed data:
   ```bash
   make db-init
   ```

4. **Ready to Go!**
   The backend API is now running and reachable at: **`http://localhost:5050`**.

---

### Method B: Running Locally (Local Development)

If you prefer to run and develop the application directly on your local system:

1. **Configure Environment Variables (`.env`):**
   Duplicate or create a `.env` file in the root directory and update its contents:
   ```env
   HTTP _ADDR=:5050
   DB_HOST=localhost
   DB_PORT=3306
   DB_USERNAME=user
   DB_PASSWORD=password123
   DB_NAME=article
   ```

2. **Run Your Local MySQL Instance:**
   Ensure MySQL is running on port `3306` with the credentials matching your `.env` file.

3. **Initialize & Seed the Database:**
   ```bash
   make db-init
   ```

4. **Run the Go Server:**
   ```bash
   go run main.go
   ```
   *Alternatively, you can use the Makefile helper:*
   ```bash
   make run
   ```

---

## 📑 API Reference (REST Endpoints)

All request and response payloads use the `application/json` format.

| Method | Endpoint | Query Parameters / Path | Description |
| :--- | :--- | :--- | :--- |
| **GET** | `/article` | `?limit=10&page=1&status=publish` | Fetches a paginated list of articles (excludes `thrash` status if no status folter is provided). |
| **GET** | `/article/{id}` | `/article/1` | Retrieves the details of a specific article by its ID. |
| **POST** | `/article` | *None* (Request Body JSON) | Creates a new article. |
| **PUT** | `/article/{id}` | `/article/1` (Request Body JSON) | Updates an`existing article's content by its ID. |
| **DELETE** | `/article/{id}` | `/article/1` | Soft deletes an article (moves status to `thrash`). |

### Paginated Response Structure (GET `/article`)
```json
{
  "data": [
    {
      "id": 1,
      "title": "Learning Golang REST API Development",
      "content": "This is an`example article content built purely for validation testing on this robust backend implementation using Go...",
      "category": "Golang",
      "created_date": "2026-07-16T17:29:27Z",
      "updated_date": "2026-07-16T17:29:27Z",
      "status": "publish"
    
  ],
  "total_rows": 25,
  "total_pages": 5
}
```

---

## 🧪 API Testing

To make testing as seamless as possible, a **Postman Collection** is included in the project root directory.

1. Open **Postman**.
2. Import the `article_api.postman_collection.json` format.
3. You can immediately run and test all scenarios (including Success states, Validation Errors, Pagination, and Soft Delete actions).