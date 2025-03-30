# syntax=docker/dockerfile:1

# Build stage
FROM golang:1.23.0 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN cd cmd/api/v1 && go build -o lumon-backend && mv ./lumon-backend /app/


# Final stage
FROM ubuntu:22.04 AS build-release-stage

# Install necessary dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Create a non-root user
RUN useradd -m nonroot
USER nonroot

# Set the working directory
WORKDIR /server

# Copy the binary and template files from the build stage
COPY --from=build-stage /app/lumon-backend  /server/

# Define a build argument for the port
ARG APP_PORT=5455

# Set the port as an environment variable
ENV APP_PORT=${APP_PORT}

# Expose the port
EXPOSE ${APP_PORT}

ENV DOCKERIZE_VERSION v0.7.0

USER root

RUN apt-get update \
    && apt-get install -y wget \
    && ARCH=$(dpkg --print-architecture) \
    && if [ "$ARCH" = "arm64" ]; then \
    wget -O /usr/local/bin/dockerize https://github.com/jwilder/dockerize/releases/download/$DOCKERIZE_VERSION/dockerize-linux-arm64-$DOCKERIZE_VERSION.tar.gz && \
    tar -xzf /usr/local/bin/dockerize -C /usr/local/bin; \
    else \
    wget -O /usr/local/bin/dockerize https://github.com/jwilder/dockerize/releases/download/$DOCKERIZE_VERSION/dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz && \
    tar -xzf /usr/local/bin/dockerize -C /usr/local/bin; \
    fi \
    && chmod +x /usr/local/bin/dockerize \
    && apt-get remove -y wget \
    && apt-get autoremove -y \
    && rm -rf /var/lib/apt/lists/*

USER nonroot

# Run the application
CMD ["dockerize", "-wait", "tcp://postgres:5432", "/server/lumon-backend"]
