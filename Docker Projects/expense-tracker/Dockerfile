FROM golang:1.21-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /expense-tracker ./cmd/server

FROM gcr.io/distroless/static-debian11
COPY --from=build /expense-tracker /expense-tracker
EXPOSE 8080
ENTRYPOINT ["/expense-tracker"]
