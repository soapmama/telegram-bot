FROM golang:1.24-alpine as build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o telegram-bot ./cmd

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /app/telegram-bot /telegram-bot

EXPOSE 4211

# Run the application
ENTRYPOINT ["/telegram-bot"] 