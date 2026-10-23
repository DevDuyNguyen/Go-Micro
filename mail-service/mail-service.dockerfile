FROM golang:1.27.0 AS build
WORKDIR /app
COPY go.sum go.mod /app/
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o mail-service ./cmd/api

FROM alpine
WORKDIR /app
COPY --from=build /app/mail-service .
EXPOSE 80
CMD ["./mail-service"]