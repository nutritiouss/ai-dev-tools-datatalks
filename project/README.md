# Forum API Wrapper & Dashboard

A full-stack application that creates a modern REST API wrapper for the ReSQL forum (https://resql.ru/forum/), with a React dashboard for browsing and searching forum content.

## Problem Description

The ReSQL forum is a popular Russian-language forum for database and SQL discussions. However, it lacks a modern REST API for programmatic access to its content. This project solves this problem by:

1. **Creating a REST API wrapper** that provides structured access to forum data (topics, posts, users, forums)
2. **Building a modern web dashboard** with search and filtering capabilities
3. **Implementing a data synchronization system** to keep the API data up-to-date with the forum

The system allows developers and users to:
- Browse forum topics and posts through a modern interface
- Search across forum content
- Access forum data programmatically via REST API
- Filter and sort topics by various criteria

## System Architecture

```
┌─────────────┐         ┌──────────────┐         ┌─────────────┐
│   Frontend  │────────▶│   Backend    │────────▶│  PostgreSQL │
│   (React)   │  HTTP   │    (Go)      │   SQL   │  Database   │
└─────────────┘         └──────────────┘         └─────────────┘
                              │
                              │ HTTP
                              ▼
                        ┌──────────────┐
                        │ ReSQL Forum  │
                        │   (Scraper)  │
                        └──────────────┘
```

### Components

- **Frontend**: React + TypeScript application with Vite
- **Backend**: Go application using Gin framework
- **Database**: PostgreSQL (production) / SQLite (development)
- **API**: RESTful API following OpenAPI 3.0 specification

## Technologies Used

### Frontend
- **React 18+**: UI framework
- **TypeScript**: Type safety
- **Vite**: Build tool and dev server
- **React Router**: Client-side routing
- **Vitest**: Testing framework

### Backend
- **Go 1.21+**: Programming language
- **Gin**: HTTP web framework
- **PostgreSQL**: Production database
- **SQLite**: Development/testing database
- **pgx**: PostgreSQL driver

### Infrastructure
- **Docker**: Containerization
- **Docker Compose**: Local development orchestration
- **GitHub Actions**: CI/CD pipeline
- **Render/Railway**: Deployment platform (optional)

## Features

- ✅ RESTful API with OpenAPI specification
- ✅ Forum topics browsing with pagination
- ✅ Topic detail view with posts
- ✅ Full-text search across topics, posts, and users
- ✅ User profiles and activity
- ✅ Responsive, modern UI
- ✅ Database migrations
- ✅ Integration tests
- ✅ Docker containerization
- ✅ CI/CD pipeline

## Project Structure

```
project/
├── backend/              # Go backend
│   ├── cmd/
│   │   └── server/       # Application entry point
│   ├── internal/
│   │   ├── api/          # HTTP handlers
│   │   ├── models/        # Data models
│   │   ├── repository/   # Database layer
│   │   ├── service/      # Business logic
│   │   └── scraper/       # Forum scraping
│   ├── tests/
│   │   └── integration/  # Integration tests
│   └── Dockerfile
├── frontend/             # React frontend
│   ├── src/
│   │   ├── components/   # React components
│   │   ├── services/      # API client
│   │   └── utils/
│   └── Dockerfile
├── api/                  # OpenAPI specifications
│   └── openapi.yaml
├── docker-compose.yml
├── .github/
│   └── workflows/
│       └── ci-cd.yml
└── README.md
```

## Setup Instructions

### Prerequisites

- Go 1.21 or higher
- Node.js 20 or higher
- Docker and Docker Compose (for containerized setup)
- PostgreSQL (for production) or SQLite (for development)

### Local Development

#### Option 1: Using Docker Compose (Recommended)

1. **Clone the repository**:
   ```bash
   git clone <repository-url>
   cd project
   ```

2. **Start all services**:
   ```bash
   docker-compose up --build
   ```

3. **Access the application**:
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080/api
   - API Health: http://localhost:8080/api/health

#### Option 2: Manual Setup

**Backend Setup**:

1. Navigate to backend directory:
   ```bash
   cd backend
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Set environment variables (optional, defaults to SQLite):
   ```bash
   export DATABASE_URL="postgres://user:password@localhost/forum_db?sslmode=disable"
   export DB_DRIVER="postgres"
   export PORT=8080
   ```

4. Run the server:
   ```bash
   go run cmd/server/main.go
   ```

**Frontend Setup**:

1. Navigate to frontend directory:
   ```bash
   cd frontend
   ```

2. Install dependencies:
   ```bash
   npm install
   ```

3. Set environment variable (optional):
   ```bash
   export VITE_API_BASE_URL=http://localhost:8080/api
   ```

4. Start development server:
   ```bash
   npm run dev
   ```

## Testing

### Backend Tests

Run unit tests:
```bash
cd backend
go test ./... -v
```

Run integration tests:
```bash
cd backend
go test ./tests/integration/... -v
```

### Frontend Tests

Run tests:
```bash
cd frontend
npm test
```

## API Documentation

The API follows the OpenAPI 3.0 specification defined in `api/openapi.yaml`.

### Endpoints

- `GET /api/health` - Health check
- `GET /api/forums` - List forums
- `GET /api/forums/:id` - Get forum by ID
- `GET /api/topics` - List topics (with filtering)
- `GET /api/topics/:id` - Get topic with posts
- `GET /api/posts` - List posts (with filtering)
- `GET /api/posts/:id` - Get post by ID
- `GET /api/users` - List users
- `GET /api/users/:id` - Get user by ID
- `GET /api/search` - Search across content

### Example Request

```bash
curl http://localhost:8080/api/topics?page=1&limit=20&sort=newest
```

## Deployment

### Using Docker

Build and run with Docker Compose:
```bash
docker-compose up --build -d
```

### Using Render

1. Connect your GitHub repository to Render
2. Create a new Web Service for the backend
3. Create a new Web Service for the frontend
4. Create a PostgreSQL database
5. Set environment variables as specified in `render.yaml`

### Environment Variables

**Backend**:
- `DATABASE_URL`: PostgreSQL connection string
- `DB_DRIVER`: Database driver (`postgres` or `sqlite3`)
- `PORT`: Server port (default: 8080)

**Frontend**:
- `VITE_API_BASE_URL`: Backend API URL

## CI/CD Pipeline

The project includes a GitHub Actions workflow (`.github/workflows/ci-cd.yml`) that:

1. Runs backend unit and integration tests
2. Runs frontend tests
3. Builds Docker images
4. Deploys to production (when tests pass)

## AI Development Tools Usage

This project was developed using AI-assisted development tools:

- **Coding Assistant**: Used for generating boilerplate code, implementing API handlers, and creating React components
- **Code Review**: AI tools helped identify potential bugs and suggest improvements
- **Documentation**: AI assisted in generating comprehensive documentation

### Development Workflow

1. **OpenAPI First**: Started with OpenAPI specification to define the API contract
2. **Frontend Development**: Built React components with centralized API client
3. **Backend Implementation**: Implemented Go backend following the OpenAPI spec
4. **Integration**: Connected frontend to backend and added tests
5. **Containerization**: Created Dockerfiles and docker-compose configuration
6. **CI/CD**: Set up GitHub Actions workflow

## Database Schema

- **forums**: Forum categories
- **topics**: Discussion topics
- **posts**: Individual posts/replies
- **users**: Forum users

See `backend/cmd/server/migrations.go` for the complete schema.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

This project is part of the AI Dev Tools Zoomcamp course.

## Future Enhancements

- [ ] Implement actual forum scraping logic
- [ ] Add authentication and authorization
- [ ] Implement real-time updates
- [ ] Add caching layer
- [ ] Implement rate limiting
- [ ] Add analytics dashboard
- [ ] Support for multiple forums
