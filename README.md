# Banners Application

Banner Service is a functional solution designed for managing and displaying banners tailored to user tags and features. Leveraging Redis for caching, PostgreSQL for data management, Docker Compose for service orchestration and Jaeger for tracing

## Prerequisites

Before proceeding, ensure you have the following installed on your system:
- Docker and Docker Compose (https://docs.docker.com/compose/install/)
- Golang (https://go.dev/doc/install)
- Goose (for database migrations) (https://github.com/pressly/goose)

## Getting Started

Follow these steps to set up and run the application:

### Step 1: Start Redis and PostgreSQL

Firstly, start the Redis and PostgreSQL services using Docker Compose. Run the following command in the root of your project directory:

```bash
docker-compose up
```

### Step 2: Apply Database Migrations

After the Redis and PostgreSQL services are up, apply the database migrations with this command:

```bash
make migration-up
```

### Step 3: Run the Application

Finally, to start the application, execute the following command:

```bash
go run ./cmd/main.go
```
