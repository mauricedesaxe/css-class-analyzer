# Use the official Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH environment variable at /go.
FROM golang:1.21.3 as builder

# Create and change to the app directory.
WORKDIR /app

# Retrieve application dependencies.
# This allows the container build to reuse cached dependencies.
COPY go.* ./
RUN go mod download

# Copy local code to the container image.
COPY . ./

# Build the binary.
# -o myapp specifies the output name of the binary.
RUN CGO_ENABLED=0 GOOS=linux go build -v -o myapp

# Use the official Debian slim image for a lean production container.
# https://hub.docker.com/_/debian
FROM debian:buster-slim
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/myapp /myapp

# Service must listen to $PORT environment variable.
# This default value facilitates local development.
ENV PORT 3000

# Run the web service on container startup.
CMD ["/myapp"]