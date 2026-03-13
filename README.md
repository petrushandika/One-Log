# 🚀 ULAM: Unified Log & Activity Monitor

ULAM is a high-performance, centralized observability platform designed to capture, store, and analyze technical logs and user activities across all your applications. Equipped with **Groq-powered AI Insight**, ULAM doesn't just store logs—it understands them.

---

## ✨ Key Features

- **🛡️ Centralized Logging**: Single endpoint (`/api/ingest`) for all your apps (Go, Node, PHP, Python, etc.).
- **🤖 AI Insight (Groq)**: Instant error summarization, Root Cause Analysis (RCA), and solution suggestions using Llama 3.3.
- **🔐 User Activity & Auth Tracking**: Detailed audit trails for login events (OAuth, Magic Link, Password) and user actions.
- **🎭 PII Masking**: Automatically censor sensitive data (passwords, tokens, secrets) before it hits the database.
- **📉 Error Grouping**: Intelligent log fingerprinting to reduce dashboard noise by grouping identical issues.
- **📩 Real-time Alerting**: Automated email notifications for `ERROR` and `CRITICAL` levels with smart throttling.
- **🧹 Auto Retention**: Automated cleanup policy to keep your database lean and performant.

---

## 🛠️ Tech Stack

| Layer | Technology |
| :--- | :--- |
| **Backend** | Go 1.26 (Gin, GORM v2) |
| **Frontend** | React 19 (Vite 7, Tailwind v4) |
| **Database** | PostgreSQL 17 (JSONB, GIN Index) |
| **AI Engine** | Groq API (Llama 3.3) |
| **Auth** | JWT via httpOnly Cookies + API Key |
| **Infra** | Docker & Docker Compose v2 |

---

## 📂 Project Navigation

The project is highly documented to ensure scalability and ease of integration.

- [**📋 PRD**](Documentation/PRD.md) - Product requirements and business logic.
- [**🏗️ Architecture**](Documentation/ARCHITECTURE.md) - Detailed system design and structure.
- [**🗄️ Data Schema**](Documentation/DATA_SCHEMA.md) - Database models and JSONB contracts.
- [**🔗 API Reference**](Documentation/API_REFERENCE.md) - Endpoints and integration guide.
- [**🛠️ Tech Stack**](Documentation/TECH_STACK.md) - Detailed package versions and tool choices.
- [**🗺️ Roadmap**](Documentation/ROADMAP.md) - Feature implementation phases.
- [**🏆 MVP Scope**](Documentation/MVP.md) - Focus for current development.

---

## 🚀 Quick Start

### 1. Prerequisites
- Docker & Docker Compose
- Groq API Key (for AI features)
- Gmail App Password (for email alerts)

### 2. Setup Environment
```bash
cp .env.example .env
# Fill in DATABASE_URL, JWT_SECRET, GROQ_API_KEY, and SMTP credentials
```

### 🏃‍♂️ Quick Start (Method 1: Makefile - RECOMMENDED)
We provide a `Makefile` for a clean and automated workflow:

```bash
# 1. Setup all dependencies (Go & NPM)
make setup

# 2. Run Database Migrations
make migrate

# 3. Build both Frontend & Backend
make build

# 4. Start Development Mode
make dev
```

### 🏃‍♂️ Quick Start (Method 2: Docker)
If you prefer Docker:
```bash
docker-compose up -d --build
```
The dashboard will be available at `http://localhost:5173` and the API at `http://localhost:8080`.

---

## 🧪 Integration Example (Go)

```go
// Send log to ULAM
ulam.Send(LogPayload{
    Category: "SYSTEM_ERROR",
    Level:    "CRITICAL",
    Message:  "Database connection refused",
    Context: map[string]interface{}{"user_id": "usr_123"},
})
```

---

## 📄 License
This project is proprietary. Created by **Petrus Handika**.
