# Build stage
FROM node:20-alpine AS builder

WORKDIR /app

# Copy frontend files
COPY frontend/package*.json ./frontend/
RUN cd frontend && npm ci

COPY frontend/ ./frontend/
RUN cd frontend && npm run build

# Production stage
FROM node:20-alpine

WORKDIR /app

# Copy backend files
COPY backend/package*.json ./
RUN npm ci --only=production

# Copy backend source
COPY backend/ ./

# Copy built frontend
COPY --from=builder /app/frontend/dist ./public

# Expose port
EXPOSE 5000

# Start server
CMD ["node", "server.js"]







