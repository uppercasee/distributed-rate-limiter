FROM golang:1.24-alpine

ENV GO111MODULE=on

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o server ./cmd/server

EXPOSE 50051

CMD ["./server"]
