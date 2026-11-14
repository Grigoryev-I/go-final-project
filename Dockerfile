FROM golang:1.26.1 AS builder

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download

COPY main.go ./
COPY pkg ./pkg

RUN CGO_ENABLED=0 GOOS=linux go build -o /build/main .

FROM alpine:latest

WORKDIR /app
COPY --from=builder /build/main ./main
COPY web ./web
RUN mkdir -p /app/data

ENV TODO_PORT=7540
ENV TODO_DBFILE=/app/data/scheduler.db

CMD ["./main"]