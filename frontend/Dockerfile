# Build stage
FROM node:20-alpine AS builder

# Set working directory
WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies (including dev dependencies for build)
RUN npm ci

# Copy source code
COPY . .

# Build the application
RUN npm run build

# Clean up node_modules to reduce image size
RUN rm -rf node_modules

# Production stage
FROM nginx:alpine

# Copy built assets from builder stage
COPY --from=builder /app/dist /usr/share/nginx/html

# Copy nginx configuration
COPY nginx.conf /etc/nginx/nginx.conf

# Create a script to inject environment variables
RUN echo '#!/bin/sh' > /docker-entrypoint.d/30-env-injection.sh && \
    echo 'echo "window.env = {" > /usr/share/nginx/html/env.js' >> /docker-entrypoint.d/30-env-injection.sh && \
    echo 'echo "  REACT_APP_ENV: \"$REACT_APP_ENV\"," >> /usr/share/nginx/html/env.js' >> /docker-entrypoint.d/30-env-injection.sh && \
    echo 'echo "  REACT_APP_API_BASE_URL: \"$REACT_APP_API_BASE_URL\"" >> /usr/share/nginx/html/env.js' >> /docker-entrypoint.d/30-env-injection.sh && \
    echo 'echo "};" >> /usr/share/nginx/html/env.js' >> /docker-entrypoint.d/30-env-injection.sh && \
    chmod +x /docker-entrypoint.d/30-env-injection.sh

# Expose port
EXPOSE 80

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:80/ || exit 1

# Start nginx
CMD ["nginx", "-g", "daemon off;"]