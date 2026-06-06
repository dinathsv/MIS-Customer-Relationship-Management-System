FROM golang:1.21-alpine

WORKDIR /app

# Initialize module and fetch PostgreSQL driver dependencies
RUN go mod init complain-backend && go get github.com/lib/pq

COPY . .

RUN go build -o main .

EXPOSE 8080

CMD ["./main"]