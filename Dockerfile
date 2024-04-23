# Base Image
FROM golang:1.22.2 as base

# Working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY ./src ./src
COPY ./template ./template

# Build the application
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o server ./src/cmd/main.go

# Create master image
FROM alpine AS master

RUN apk add --no-cache libc6-compat 

# Working directory
WORKDIR /app

# Copy execute file
COPY --from=base /app/server ./

# Set ENV to production
ENV GO_ENV production

# Expose port 8080
EXPOSE 8080

# Run the application
CMD ["./server"]