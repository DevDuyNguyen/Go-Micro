FROM golang:1.26.4-alpine AS builder
WORKDIR /app
COPY . /app
RUN CGO_ENABLED=0 go build -o brokerApp ./cmd/api
RUN chmod +x brokerApp

FROM alpine:latest
RUN mkdir /app
COPY --from=builder /app/brokerApp /app
EXPOSE 80
CMD [ "/app/brokerApp" ]