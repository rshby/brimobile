FROM golang:1.21.1-alpine as builder

WORKDIR /app

COPY ./ ./

RUN go mod tidy
RUN go build -o ./bin/brimobile ./server.go

FROM alpine:3

WORKDIR /app

COPY --from=builder /app/.env ./
COPY --from=builder /app/bin/brimobile ./

EXPOSE 7130
CMD ./brimobile