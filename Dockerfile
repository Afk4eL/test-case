# syntax=docker/dockerfile:1
FROM golang:1.22.5

WORKDIR /test-case

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /server ./cmd/test-case

EXPOSE 8080

CMD ["/server", "/test-case/config/prod.yaml"]