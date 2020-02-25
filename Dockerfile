FROM golang:latest

WORKDIR /api

COPY go.mod ./

RUN go mod download

COPY . .

RUN go build -o main .

EXPOSE 8080

CMD ["go run main.go"]