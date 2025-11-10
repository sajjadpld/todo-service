FROM golang:1.24-alpine3.20 AS builder

WORKDIR /app

COPY . .

ENV PATH=$PATH:/go/bin

RUN apk update && apk add make g++ bash

RUN go install github.com/swaggo/swag/cmd/swag@latest

RUN make swag

ENV GOPROXY https://goproxy.io

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main /app/cmd/main.go

FROM alpine:3.20 AS production

WORKDIR /app

RUN mkdir -p /app/locale /app/logs /app/schema/psql && chmod 755 /app/locale /app/logs /app/schema/psql

COPY schema/psql /app/schema/psql

COPY internal/adapter/locale/translation /app/locale

COPY .env.example .

RUN mv .env.example .env

COPY --from=builder /app/main .

RUN chmod +x main

EXPOSE 8080

CMD ["./main"]