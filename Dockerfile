FROM golang:1.26

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN go build -o myapp ./cmd

EXPOSE 8000

CMD ["./myapp"]
