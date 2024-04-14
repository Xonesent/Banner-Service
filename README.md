# Banners Application

Banner Service is a functional solution designed for managing and displaying banners tailored to user tags and features. Leveraging Redis for caching, PostgreSQL for data management, Docker Compose for service orchestration and Jaeger for tracing

## Prerequisites

Before proceeding, ensure you have the following installed on your system:
- Docker and Docker Compose (https://docs.docker.com/compose/install/)
- Golang (https://go.dev/doc/install)
- Goose (for database migrations) (https://github.com/pressly/goose)

## Getting Started

Follow these steps to set up and run the application:

### Step 1: Start Redis, PostgreSQL and Jaeger

Firstly, start the Redis, PostgreSQL and Jaeger services using Docker Compose. Run the following command in the root of your project directory:

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

## My impressions and thoughts

Throughout the project I have left readme files that will emphasize my reasoning and also leave questions that I would like to get answers to in the feedback. I'm sure some of the questions will be debatable, so I won't be too upset if there is no answer to them  
- In internal = about performance and functionality
- In migrations = about constraints and why I made these particular table structures.
- In pkg = about architecture, my assumptions about what could be done better

Would be happy to get feedback / advice!  
Email - 1pyankov.d.s@gmail.com  
Telegram - @Xonesent