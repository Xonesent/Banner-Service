FROM golang:latest AS builder

WORKDIR .

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o banners cmd/main.go

EXPOSE 8080

CMD ["sh", "-c", "sleep 5 && ./banners"]