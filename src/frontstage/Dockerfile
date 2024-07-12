# Stage 1: Build the application
FROM node:20-alpine AS builder

# Create and change to the app directory
WORKDIR /app

# Install pnpm
RUN npm install -g pnpm

# Copy package.json and pnpm-lock.yaml
COPY package.json pnpm-lock.yaml ./

# Install dependencies
RUN pnpm install

# Copy all files
COPY . .

# Build the application
RUN pnpm run build

# Stage 2: Run the application
FROM node:20-alpine

# Create and change to the app directory
WORKDIR /app

# Install pnpm
RUN npm install -g pnpm

# Copy only necessary files from the builder stage
COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app/.next ./.next
COPY --from=builder /app/public ./public
COPY --from=builder /app/package.json ./

COPY .env.local ./

# Expose the port
EXPOSE 3000

# Start the application
CMD ["pnpm", "start"]
