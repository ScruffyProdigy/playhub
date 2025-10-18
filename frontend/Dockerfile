# Stage 1: build the React app
FROM node:24-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

# Stage 2: serve it with nginx
FROM nginx:latest

# Copy built static assets
COPY --from=builder /app/dist /usr/share/nginx/html

# Copy startup scripts into the official nginx entrypoint directory
COPY docker-entrypoint.d/ /docker-entrypoint.d/

# Make sure scripts are executable
RUN chmod +x /docker-entrypoint.d/*.sh

# (Optional) Expose port
EXPOSE 80

# Use the default nginx entrypoint which automatically runs scripts in /docker-entrypoint.d
ENTRYPOINT ["/docker-entrypoint.sh"]
CMD ["nginx", "-g", "daemon off;"]
