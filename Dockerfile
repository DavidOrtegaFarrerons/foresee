FROM golang:1.24

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

CMD ["go", "run", "./cmd/web"]