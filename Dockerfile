# Build stage
FROM node:20-alpine AS builder

# Install git for submodule operations
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy package files
COPY package*.json ./
COPY tsconfig.json ./

# Install dependencies
RUN npm install

# Copy source code
COPY . .

# Initialize git and submodules if they exist
RUN if [ -f .gitmodules ]; then \
        git init && \
        git submodule update --init --recursive; \
    fi

# Set environment variables for Next.js
ENV NEXT_TELEMETRY_DISABLED=1
ENV NODE_ENV=production
ENV NEXT_SHARP_PATH=/app/node_modules/sharp
ENV NEXT_RUNTIME=edge

# Build the application
RUN npm run build

# Production stage
FROM node:20-alpine

WORKDIR /app

# Set environment variables
ENV NODE_ENV=production
ENV NEXT_TELEMETRY_DISABLED=1
ENV NEXT_SHARP_PATH=/app/node_modules/sharp
ENV NEXT_RUNTIME=edge

# Copy built assets from builder
COPY --from=builder /app/.next ./.next
COPY --from=builder /app/package*.json ./
COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app/next.config.js ./

# Create and copy public directory
RUN mkdir -p public
COPY --from=builder /app/public ./public

# Copy submodule contents if they exist
COPY --from=builder /app/jaxstats ./jaxstats

# Copy requirements and install dependencies
COPY requirements.txt .
RUN pip install -r requirements.txt

# Copy backend code
COPY . .

# Copy frontend build into static directory
COPY frontend/build ./app/static/

# Create data directory for replays
RUN mkdir -p data/replays

# Expose the application port
EXPOSE 3000

# Start the application
CMD ["npm", "start"] 