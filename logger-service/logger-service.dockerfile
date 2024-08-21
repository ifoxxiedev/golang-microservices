FROM golang:1.23-alpine as builder

RUN mkdir /app
COPY . /app
WORKDIR /app

# RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o loggerApp ./cmd/api
RUN CGO_ENABLED=0 go build -o loggerApp ./cmd/api
RUN chmod +x /app/loggerApp

FROM alpine:latest
RUN mkdir /app

COPY --from=builder /app/loggerApp /app

CMD ["./app/loggerApp"]