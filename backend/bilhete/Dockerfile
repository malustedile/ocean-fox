FROM golang:1.22-alpine

WORKDIR /app

COPY . .

RUN go mod tidy && go build -o bilhete main.go

EXPOSE 3002

CMD ["./bilhete"]