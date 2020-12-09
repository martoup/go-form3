FROM golang:1.14.3-alpine
WORKDIR /app
COPY go.mod .
COPY test ./test
COPY form3 ./form3
ENTRYPOINT ./test/wait.sh
