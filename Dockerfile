FROM golang:1.24-alpine AS build

WORKDIR /app

RUN apk add --no-cache make

# Copy and download dependencies
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the Go application
RUN make build-binary-in-docker

# Stage 2: Final stage
FROM alpine:edge

# Set the working directory
WORKDIR /app

# Copy the binary from the build stage
COPY --from=build /app/build/gorate .

# Set the timezone and install CA certificates
RUN apk --no-cache add ca-certificates

# Set the entrypoint command
ENTRYPOINT ["/app/gorate"]