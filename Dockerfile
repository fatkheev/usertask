FROM golang:1.22-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /app/usertask cmd/main.go

RUN chmod +x /app/usertask

EXPOSE 8080

CMD ["/app/usertask"]