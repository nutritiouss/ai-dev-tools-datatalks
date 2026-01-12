# AI Development Guidelines

This document provides guidelines for AI assistants working on this project.

## Backend (Go)

### Dependency Management
- Use Go modules (`go mod`) for dependency management
- Add dependencies with: `go get <package>`
- Update dependencies with: `go get -u <package>`
- Run `go mod tidy` to clean up unused dependencies

### Project Structure
- Main entry point: `backend/cmd/server/main.go`
- Internal packages: `backend/internal/`
- API handlers: `backend/internal/api/`
- Business logic: `backend/internal/service/`
- Database access: `backend/internal/repository/`
- Data models: `backend/internal/models/`
- Scraper logic: `backend/internal/scraper/`

### Code Style
- Follow Go standard formatting (`gofmt`)
- Use meaningful variable and function names
- Add comments for exported functions and types
- Keep functions focused and small
- Use interfaces for testability

### Testing
- Unit tests: `*_test.go` files alongside source files
- Integration tests: `backend/tests/integration/`
- Run tests with: `go test ./...`
- Use table-driven tests when appropriate

### Database
- Use `github.com/jackc/pgx/v5` for PostgreSQL
- Use `github.com/golang-migrate/migrate` for migrations
- Support both PostgreSQL (production) and SQLite (development/testing)
- Use connection pooling
- Handle database errors gracefully

### API Framework
- Use Gin framework for HTTP routing
- Follow RESTful conventions
- Return proper HTTP status codes
- Use middleware for common concerns (logging, CORS, auth)

## Frontend (React + TypeScript)

### Dependency Management
- Use npm or yarn for package management
- Add dependencies with: `npm install <package>` or `yarn add <package>`
- Use Vite as the build tool

### Project Structure
- Components: `frontend/src/components/`
- API client: `frontend/src/services/api.ts` (centralized)
- Utilities: `frontend/src/utils/`
- Main app: `frontend/src/App.tsx`

### Code Style
- Use TypeScript for type safety
- Use functional components with hooks
- Keep components small and focused
- Extract reusable logic into custom hooks
- Use meaningful component and variable names

### Testing
- Use Jest and React Testing Library
- Test user interactions, not implementation details
- Write tests for core business logic
- Run tests with `npm test` or `yarn test`

### API Communication
- All API calls must go through `src/services/api.ts`
- Use async/await for API calls
- Handle errors appropriately
- Show loading states

## General Guidelines

### Git
- Write clear commit messages
- Commit frequently with logical units of work
- Don't commit generated files or dependencies

### Documentation
- Update README.md when adding features
- Document complex logic with comments
- Keep API documentation (OpenAPI) up to date

### Error Handling
- Handle errors gracefully
- Provide meaningful error messages
- Log errors appropriately
- Don't expose sensitive information in errors
