# CareerCompass Backend

A Go-based REST API backend built with Gin framework and PostgreSQL database.

## 📋 Table of Contents

- [Prerequisites](#prerequisites)
- [Tech Stack](#tech-stack)
- [Project Structure](#project-structure)
- [Getting Started](#getting-started)
- [Environment Variables](#environment-variables)
- [Database Management](#database-management)
- [Running the Application](#running-the-application)
- [API Endpoints](#api-endpoints)
- [Troubleshooting](#troubleshooting)

## 🔧 Prerequisites

Before you begin, ensure you have the following installed on your system:

- **Go** (version 1.25.6 or higher) - [Download Go](https://golang.org/dl/)
- **Docker** and **Docker Compose** (version 29.1.5) - [Install Docker](https://docs.docker.com/get-docker/)
- **Git** - [Install Git](https://git-scm.com/downloads)

## 🛠 Tech Stack

- **Framework**: [Gin](https://gin-gonic.com/) - HTTP web framework
- **Database**: PostgreSQL 16
- **Database Driver**: [pgx/v5](https://github.com/jackc/pgx) - PostgreSQL driver and toolkit
- **Migrations**: [golang-migrate](https://github.com/golang-migrate/migrate) - Database migration tool
- **Environment**: [godotenv](https://github.com/joho/godotenv) - Environment variable management
- **CORS**: [gin-contrib/cors](https://github.com/gin-contrib/cors) - CORS middleware
- **Authentication**: [jwt](https://github.com/golang-jwt/jwt) - JSON Web Tokens
- **Security**: [bcrypt](https://golang.org/x/crypto/bcrypt) - Password hashing

## 📁 Project Structure

```
careercompass-backend/
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/                  # Configuration management
│   ├── database/                # Database connection setup
│   ├── handlers/                # HTTP request handlers
│   ├── models/                  # Data models
│   ├── router/                  # Route definitions
│   └── services/                # Business logic & Authentication
├── migrations/                  # Database migration files
│   ├── 000001_create_users_table.up.sql
│   └── 000001_create_users_table.down.sql
├── docker-compose.yml           # Docker services configuration
├── .env.example                 # Example environment variables
├── go.mod                       # Go module dependencies
└── go.sum                       # Go module checksums
```

## 🚀 Getting Started

Follow these steps to set up the project from scratch:

### 1. Clone the Repository

```bash
git clone <repository-url>
cd careercompass-backend
```

### 2. Set Up Environment Variables

Copy the example environment file and configure it:

```bash
cp .env.example .env
```

The `.env` file should contain:

```env
# Server Configuration
PORT=4546
GIN_MODE=debug

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=appuser
DB_PASSWORD=secretpassword
DB_NAME=careerdb
DB_SSLMODE=disable

# CORS Configuration
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173

# JWT Configuration
JWT_SECRET=your-secret-key-change-this-in-production
JWT_EXPIRATION=24h
```

> **Note**: Modify these values according to your environment needs.

### 3. Start Docker Services

Start the PostgreSQL database and Adminer (database management tool):

```bash
docker compose up -d
```

This will start:
- **PostgreSQL** on port `5432`
- **Adminer** (database UI) on port `8090`

### 4. Install Go Dependencies

```bash
go mod download
```

### 5. Run the Application

```bash
go run cmd/api/main.go
```

The server will start on `http://localhost:4546`

## 🔐 Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `4546` |
| `GIN_MODE` | Gin mode (debug/release) | `debug` |
| `DB_HOST` | Database host | `localhost` |
| `DB_PORT` | Database port | `5432` |
| `DB_USER` | Database username | `appuser` |
| `DB_PASSWORD` | Database password | `secretpassword` |
| `DB_NAME` | Database name | `careerdb` |
| `DB_SSLMODE` | SSL mode for database | `disable` |
| `CORS_ALLOWED_ORIGINS` | Allowed CORS origins | `http://localhost:3000,http://localhost:5173` |
| `JWT_SECRET` | Secret key for JWT signing | `your-secret-key...` |
| `JWT_EXPIRATION` | Token expiration time | `24h` |

## 🗄 Database Management

### Accessing Adminer

Adminer is a web-based database management tool. Access it at:

```
http://localhost:8090
```

**Login credentials:**
- System: `PostgreSQL`
- Server: `db`
- Username: `appuser`
- Password: `secretpassword`
- Database: `careerdb`

### Database Migrations

Migrations are automatically run when the application starts. They are located in the `migrations/` directory.

**Migration file naming convention:**
```
{version}_{description}.up.sql    # For applying migrations
{version}_{description}.down.sql  # For rolling back migrations
```

**Creating new migrations:**

1. Create two files in the `migrations/` directory:
   ```
   000002_your_migration_name.up.sql
   000002_your_migration_name.down.sql
   ```

2. Add your SQL in the `.up.sql` file (for applying changes)
3. Add rollback SQL in the `.down.sql` file (for reverting changes)

### Manual Migration Commands

If you need to run migrations manually using the `migrate` CLI:

```bash
# Install migrate CLI
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Run migrations up
migrate -path migrations -database "postgresql://appuser:secretpassword@localhost:5432/careerdb?sslmode=disable" up

# Rollback last migration
migrate -path migrations -database "postgresql://appuser:secretpassword@localhost:5432/careerdb?sslmode=disable" down 1
```

## 🏃 Running the Application

### Development Mode

```bash
go run cmd/api/main.go
```

### Build and Run

```bash
# Build the binary
go build -o bin/api cmd/api/main.go

# Run the binary
./bin/api
```

### Using Docker (Optional)

If you want to containerize the entire application:

```bash
docker compose up -d --build
```

## 🌐 API Endpoints

### Health Check

```http
GET /health
```

**Response:**
```json
{
  "status": "ok",
  "timestamp": "2026-01-24T20:16:11+07:00"
}
```

### Auth Endpoints

#### Register

```http
POST /api/auth/register
```

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "securepassword123",
  "confirm_password": "securepassword123",
  "display_name": "John Doe",
  "gender": "male"
}
```

**Response:**
```json
{
  "user": {
    "id": "uuid-string",
    "email": "user@example.com",
    "display_name": "John Doe",
    "gender": "male",
    "created_at": "timestamp",
    "updated_at": "timestamp"
  },
  "token": "jwt-token-string"
}
```

#### Login

```http
POST /api/auth/login
```

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "securepassword123"
}
```

**Response:**
```json
{
  "user": { ... },
  "token": "jwt-token-string"
}
```

### User Endpoints

#### Get Users

```http
GET /api/users
```

**Response:** Returns a list of users.

## 🐛 Troubleshooting

### Port Already in Use

If port `4546` is already in use:

```bash
# Check what's using the port
lsof -i :4546

# Kill the process (replace PID with actual process ID)
kill -9 <PID>
```

### Database Connection Issues

1. Ensure Docker containers are running:
   ```bash
   docker compose ps
   ```

2. Check database logs:
   ```bash
   docker compose logs db
   ```

3. Verify database credentials in `.env` match `docker-compose.yml`

### Migration Errors

If migrations fail:

1. Check migration files for syntax errors
2. Manually connect to the database and verify the schema
3. Check the `schema_migrations` table to see which migrations have been applied

### Docker Issues

Reset Docker services:

```bash
# Stop all services
docker compose down

# Remove volumes (WARNING: This deletes all data)
docker compose down -v

# Restart services
docker compose up -d
```

## 📝 Development Workflow

1. **Make changes** to your code
2. **Run the application** to test changes
3. **Create migrations** if database schema changes are needed
4. **Test endpoints** using tools like Postman or curl
5. **Commit changes** to version control

## 🤝 Contributing

1. Create a new branch for your feature
2. Make your changes
3. Test thoroughly
4. Submit a pull request

## 📄 License

[Add your license information here]

---

**Happy Coding! 🚀**
