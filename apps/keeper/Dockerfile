# syntax=docker/dockerfile:1

FROM golang:1.21.3-alpine

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY apps/keeper/go.mod ./
COPY apps/keeper/go.sum ./
RUN go mod download

COPY ./apps/keeper .

RUN CGO_ENABLED=0 GOOS=linux go build -o /keeper-build

EXPOSE 1323

# Run
CMD ["/keeper-build"]
