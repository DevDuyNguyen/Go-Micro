#builder
FROM golang:1.27.0 AS build
WORKDIR /app
COPY go.mod go.sum /app/
RUN go mod download
COPY . /app/
RUN CGO_ENABLED=0 go build -o authentication_service ./cmd/api/

FROM alpine
WORKDIR /app
COPY --from=build /app/ /app/
EXPOSE 80
CMD ["/app/authentication_service"]
