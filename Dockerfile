FROM golang:1.20-alpine

# Install bash
RUN apk add --no-cache bash

# setup go build
COPY go.mod go.sum /build/
WORKDIR /build
RUN go mod download
COPY paser.go twic210-874.pgn /build/
RUN go build -o main .

#setup wait-for-it
COPY wait-for-it.sh /usr/local/bin/wait-for-it
RUN chmod +x /usr/local/bin/wait-for-it

EXPOSE 8080

CMD ["/usr/local/bin/wait-for-it", "mysql:3306", "--timeout=30", "--strict", "--", "./main"]