# BTNG Sovereign Gold Standard - Production Dockerfile
# Multi-stage build with Nginx reverse proxy

# ---------- Stage 1: Build the Next.js API (port 3003) ----------
FROM node:22.11.0-alpine3.21 AS api-builder
WORKDIR /app

# Copy package files and install dependencies
COPY btng-api/package*.json ./btng-api/
RUN cd btng-api && npm ci --only=production

# Copy built .next folder and other files
COPY btng-api/.next ./btng-api/.next
COPY btng-api/public ./btng-api/public
COPY btng-api/package.json ./btng-api/

# ---------- Stage 2: Build the Genesis JAR (port 8080) ----------
FROM eclipse-temurin:21.0.8_9-jdk-alpine AS genesis-builder
WORKDIR /app

# Copy pre-built JAR
COPY genesis-app/build/libs/*.jar ./genesis.jar

# ---------- Stage 3: Nginx configuration ----------
FROM alpine:3.21.3 AS nginx-builder

# Install nginx and create configuration
RUN apk add --no-cache nginx wget

# Create nginx configuration directory
RUN mkdir -p /etc/nginx/http.d /var/log/nginx /var/cache/nginx

# Create BTNG nginx configuration
RUN echo 'server { \
    listen 80; \
    server_name localhost; \
    \
    # API routes - proxy to Next.js API \
    location /api/ { \
        proxy_pass http://localhost:3003; \
        proxy_set_header Host $host; \
        proxy_set_header X-Real-IP $remote_addr; \
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for; \
        proxy_set_header X-Forwarded-Proto $scheme; \
    } \
    \
    # Genesis platform routes \
    location /genesis/ { \
        proxy_pass http://localhost:8080/; \
        proxy_set_header Host $host; \
        proxy_set_header X-Real-IP $remote_addr; \
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for; \
        proxy_set_header X-Forwarded-Proto $scheme; \
    } \
    \
    # Health check \
    location /health { \
        access_log off; \
        return 200 "healthy\n"; \
        add_header Content-Type text/plain; \
    } \
    \
    # Default location - serve API \
    location / { \
        proxy_pass http://localhost:3003; \
        proxy_set_header Host $host; \
        proxy_set_header X-Real-IP $remote_addr; \
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for; \
        proxy_set_header X-Forwarded-Proto $scheme; \
    } \
}' > /etc/nginx/http.d/btng.conf

# Remove default nginx config
RUN rm -f /etc/nginx/http.d/default.conf

# ---------- Stage 4: Runtime image ----------
FROM eclipse-temurin:21.0.8_9-jre-alpine AS runtime
WORKDIR /app

# Install Node.js and Nginx in the runtime image
RUN apk add --no-cache nodejs npm nginx wget

# Copy built assets from previous stages
COPY --from=api-builder /app/btng-api ./btng-api
COPY --from=genesis-builder /app/genesis.jar ./genesis.jar
COPY --from=nginx-builder /etc/nginx/http.d/btng.conf /etc/nginx/http.d/btng.conf

# Create necessary directories for nginx
RUN mkdir -p /var/log/nginx /var/cache/nginx /run/nginx

# Create a non-root user
RUN addgroup -g 1001 -S nodejs && \
    adduser -S nextjs -u 1001 && \
    addgroup nextjs nginx

# Change ownership of the app directory
RUN chown -R nextjs:nodejs /app && \
    chown -R nginx:nginx /var/log/nginx /var/cache/nginx /run/nginx

# Expose port 80 for nginx
EXPOSE 80

# Health check - check nginx health endpoint
HEALTHCHECK --interval=30s --timeout=10s --start-period=60s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost/health || exit 1

# Entrypoint - start nginx and both services
COPY docker-entrypoint.sh /docker-entrypoint.sh
RUN chmod +x /docker-entrypoint.sh
ENTRYPOINT ["/docker-entrypoint.sh"]