# ⚙️ Backend Documentation

## Overview
ULAM Backend is built using **Golang 1.26** with the **Gin** web framework and **GORM** for ORM. It follows a clean architecture pattern to ensure maintainability and testability.

## Project Structure
```text
Backend/
├── cmd/
│   └── api/                # Application entry point
├── internal/
│   ├── domain/             # Entities and models
│   ├── handler/            # Delivery layer (HTTP handlers)
│   ├── service/            # Business logic layer
│   ├── repository/         # Data access layer
│   └── middleware/         # Gin middlewares (Auth, Logging, etc.)
├── pkg/
│   ├── database/           # DB connection & migration
│   ├── ai/                 # Groq API integration
│   ├── notifications/      # Email/SMTP service
│   └── utils/              # Shared helpers
├── .env.example            # Environment template
├── go.mod                  # Dependencies
└── README.md               # Ini file utama
├── docs/                   # Dokumentasi teknis spesifik backend
│   ├── ARCHITECTURE_DETAILS.md
│   └── SETUP_GUIDE.md
```

## Setup & Installations
Untuk panduan instalasi mendalam, silakan baca [docs/SETUP_GUIDE.md](docs/SETUP_GUIDE.md).

### 1. Prerequisites
- **Go 1.26+**
- **PostgreSQL 17**
- **Groq API Key** (for AI Analysis)

### 2. Environment Configuration
Copy the `.env.example` file to `.env`:
```bash
cp .env.example .env
```
Fill in the required values:
- `DATABASE_URL`: Your PostgreSQL connection string.
- `JWT_SECRET`: A random string for token signing.
- `SMTP_*`: Credentials for email notifications.
- `GROQ_API_KEY`: Your API key from Groq Console.

### 3. Install Dependencies
```bash
go mod tidy
```

## Running the Application

### Development Mode
```bash
go run cmd/api/main.go
```
The server will start at `http://localhost:8080`.

### Build for Production
```bash
go build -o ulam-api cmd/api/main.go
./ulam-api
```

## API Endpoints
See [Documentation/API_REFERENCE.md](../Documentation/API_REFERENCE.md) for full details.
- `POST /api/ingest`: Log ingestion (API Key auth).
## Docker & Deployment
This project includes a `Dockerfile` for containerized deployment.
To run via Docker (from root):
```bash
docker-compose up -d --build backend
```

## CI/CD
Automated builds are handled via `.github/workflows/deployment.yml`.
