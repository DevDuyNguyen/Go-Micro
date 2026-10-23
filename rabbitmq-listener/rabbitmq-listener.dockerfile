FROM golang:1.27.0 AS build
WORKDIR /app
COPY go.mod go.sum /app/
RUN go mod download
COPY . /app/
RUN CGO_ENABLED=0 go build -o rabbitmq-listener /app/cmd/

FROM alpine
WORKDIR /app
COPY --from=build /app/rabbitmq-listener /app/
CMD [ "/app/rabbitmq-listener" ]