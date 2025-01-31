# Book Microservice

## Overview
The **Book** Microservice is part of the **Library Management System**.

## Microservices

- **UserService/AuthService**: Handles user authentication and authorization.
- **BookService**: Manages books and stock.
- **CategoryService**: Manages book categories.
- **AuthorService**: Manages authors.

This microservice is built using:
- **Golang** for the backend.
- **PostgreSQL** for the database.
- **Docker** for containerization and deployment.
- **Docker Hub** for storing Docker images.

---

## **Technologies Used**
- **Programming Language**: Golang.
- **Database**: PostgreSQL.
- **Communication**: gRPC.
- **Middleware**: JWT for authentication.
- **Containerization**: Docker & Docker Compose.

---

## **API Documentation**
### REST API Endpoints
| HTTP Method | Endpoint                      | Description                     |
|-------------|-------------------------------|---------------------------------|
| `GET`       | `/api/v1/books`               | Get all books and search books  |
| `POST`      | `/api/v1/books`               | Create a new books              |
| `GET`       | `/api/v1/books/:id`           | Get details of a specific books |
| `PUT`       | `/api/v1/books/:id`           | Update a specific books         |
| `DELETE`    | `/api/v1/books/:id`           | Delete a specific books         |
| `GET`       | `/api/v1/books/recomendation` | Get recomendation books for user|

### gRPC Endpoints
| RPC Method          | Description                     |
|---------------------|---------------------------------|
| `DecreaseStock`     | Decrease the stock of a book    |
| `IncreaseStock`     | Increase the stock of a book    |
---

## Installation

### Prerequisites
- Install [Go](https://go.dev/doc/install)
- Install [PostgreSQL](https://www.postgresql.org/download/)
- Install [Docker](https://docs.docker.com/get-docker/)
- Install [gRPC](https://grpc.io/docs/languages/go/quickstart/)

### Running Without Docker

1. Clone the repository:
   ```sh
   git clone https://github.com/naufalhakm/library-api-book.git
   cd library-api-book
   ```
2. Setup environment variables (.env file):
   ```sh
   DB_HOST=localhost
   DB_PORT=5432
   DB_USERNAME=user
   DB_PASSWORD=password
   DB_DATABASE=library
   ```
3. Run PostgreSQL locally.
4. Start book microservice:
   ```sh
   go run cmd/server/main.go
   ```

### Running With Docker

1. Build and run services:
   ```sh
   docker-compose up -d
   ```

---

### Live Server

The microservice is running at:
http://35.240.139.186:8081/