FROM golang:1.19

WORKDIR /app

# Copy everything from the current directory to the PWD (Present Working Directory) inside the container
COPY . .

# Download all the dependencies
RUN go mod download

# Install the package
RUN go build -o ./cmd/app/app ./cmd/app/...

# This container exposes port 8080 to the outside world
EXPOSE 8080

# Run the executable
CMD ["./cmd/app/app"]