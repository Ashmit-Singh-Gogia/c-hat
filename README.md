<div align="center">
# 💬 c-hat
</div>

---

## 📖 What it is 

> **For anyone who's never touched a line of code — this section is for you.**
A simple chat app with features like Google oauth and direct messaging available.
Future scopes mentioned below
---

## ✨ Key Features

- 🔐 **Google OAuth Login / Logout** — Sign in with your existing Google account. No new passwords.
- 👤 **Automatic User Registration** — New users are seamlessly registered in the database on their first login.
- 🛡️ **JWT-Protected API Routes** — Every sensitive endpoint is guarded by a signed JSON Web Token, validated on every request.
- 💬 **Direct Messaging** — Send messages between registered users via a clean HTTP API.
- 📜 **Chat History Retrieval** — Fetch the full conversation history between any two users.
- 🌐 **CORS Configured for React** — The backend is pre-configured to accept requests from the React frontend running on `localhost:5173`.
- 🏗️ **Clean Architecture** — A layered, maintainable codebase with clear separation of concerns, ready to scale.

---

## 🛠️ Tech Stack

| Category       | Technology                                                                 |
|----------------|----------------------------------------------------------------------------|
| **Language**   | [Go (Golang)](https://go.dev/)                                             |
| **Web Framework** | [Gin](https://gin-gonic.com/)                                           |
| **Database**   | [PostgreSQL](https://www.postgresql.org/)                                  |
| **ORM**        | [GORM](https://gorm.io/)                                                   |
| **Authentication** | [Goth](https://github.com/markbates/goth) (Google OAuth 2.0) + [JWT](https://github.com/golang-jwt/jwt) |
| **Frontend**   | [React](https://react.dev/) + [Vite](https://vitejs.dev/)                  |

---

## 🚀 Getting Started

Follow these steps to get c-hat running on your local machine.

### Prerequisites

Make sure you have the following installed:

- **Go** `>= 1.21` — [Download here](https://go.dev/dl/)
- **PostgreSQL** `>= 14` — [Download here](https://www.postgresql.org/download/)
- **Node.js** `>= 18` & **npm** — [Download here](https://nodejs.org/)
- A **Google Cloud Console** project with OAuth 2.0 credentials configured.

### 1. Clone the Repository

```bash
git clone https://github.com/your-username/c-hat.git
cd c-hat
```

### 2. Configure Environment Variables

Create a `.env` file in the root of the project by copying the example below.

Then open `.env` and fill in your values:

# ── Server
PORT=8082

# ── Database 
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_postgres_user
DB_PASSWORD=your_postgres_password
DB_NAME=c_hat_db

# ── Google OAuth
# Obtain these from: https://console.cloud.google.com/
GOOGLE_CLIENT_ID=your_google_client_id
GOOGLE_CLIENT_SECRET=your_google_client_secret
GOOGLE_CALLBACK_URL=http://localhost:8080/auth/google/callback

# ── JWT 
JWT_SECRET=a_very_long_and_random_secret_string

# ── Session
SESSION_SECRET=another_random_secret_for_sessions

### 3. Set Up the Database

Create the PostgreSQL database before running the application:

```bash
psql -U your_postgres_user -c "CREATE DATABASE c_hat_db;"
```

GORM will handle all table migrations automatically when the server starts.

### 4. Run the Backend

```bash
# Install Go dependencies
go mod tidy

# Start the Go server
go run ./cmd/main.go
```

The API server will start at `http://localhost:8082`.

### 5. Run the Frontend

```bash
# Navigate to the frontend directory
cd client

# Install Node dependencies
npm install

# Start the Vite development server
npm run dev
```

The React app will be available at `http://localhost:5173`.

---

## 🏛️ Architecture

c-hat follows a **Clean (Layered) Architecture** to ensure the codebase is modular, testable, and easy to reason about. The rule is simple: **dependencies only point inward**.

```
HTTP Request
     │
     ▼
┌─────────────────────────────────────┐
│           Handlers (API Layer)       │  ← Parses HTTP requests, calls Services,
│  internal/handlers/                  │    returns HTTP responses. No business logic.
└──────────────────┬──────────────────┘
                   │ calls
                   ▼
┌─────────────────────────────────────┐
│         Services (Business Layer)    │  ← Contains all business rules and logic.
│  internal/services/                  │    Orchestrates data flow. No DB calls.
└──────────────────┬──────────────────┘
                   │ calls
                   ▼
┌─────────────────────────────────────┐
│       Repositories (Data Layer)      │  ← The only layer that talks to PostgreSQL
│  internal/repositories/             │    via GORM. Returns domain models.
└──────────────────┬──────────────────┘
                   │ reads/writes
                   ▼
           [ PostgreSQL DB ]
```

| Directory                    | Responsibility                                                         |
|------------------------------|------------------------------------------------------------------------|
| `cmd/`                       | Application entry point (`main.go`)                                    |
| `internal/handlers/`         | HTTP request/response parsing, route definitions                       |
| `internal/services/`         | Business logic, validation, orchestration                              |
| `internal/repositories/`     | All database queries via GORM                                          |
| `internal/models/`           | GORM struct definitions representing database tables                   |
| `frontend/`                    | React + Vite frontend application                                      |

---

## 🗺️ Roadmap // Future Scope

The following features are planned in order of priority:

- [ ] **🔴 Real-Time Messaging via WebSockets**
  Upgrade from the current HTTP request-response model to persistent WebSocket connections, enabling truly instant, bi-directional messaging without polling.

- [ ] **🟠 Modular OAuth Provider Interface**
  Refactor the current Google OAuth implementation into a clean `OAuthProvider` interface. This will allow adding new providers (GitHub, Apple, etc.) through a single, unified `Login()` / `Logout()` method pattern.

- [ ] **🟡 Email & Password Authentication**
  Implement traditional credential-based authentication alongside OAuth, including secure password hashing (bcrypt), email verification flow, and password reset functionality.

- [ ] **🟢 Voice & Video Calls**
  Integrate peer-to-peer Voice and Video calling using WebRTC, with signaling handled via the existing WebSocket infrastructure.

---

<div align="center">

Built with ❤️ by **Ashmit Singh Gogia**

</div>
