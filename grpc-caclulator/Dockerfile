FROM golang:latest

WORKDIR /app


COPY . .

RUN go mod download 

RUN go build -o main ./cmd/app

CMD ["sh","-c","sleep 5 && go run ./cmd/app/main.go"]