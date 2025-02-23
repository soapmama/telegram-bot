FROM golang:1.24-alpine AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /telegram-bot

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /telegram-bot /telegram-bot

EXPOSE 4211

USER nonroot:nonroot

# Build the application
RUN go build -o main .

USER nonroot:nonroot

# Run the application
ENTRYPOINT ["/telegram-bot"] 