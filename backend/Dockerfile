FROM golang:1.24-alpine

WORKDIR /app

COPY . .
RUN go mod download

RUN CGO_ENABLED=0 \
    go build \
    -ldflags "-s -w" \
    -trimpath \
    -o backend \
    .

EXPOSE 8080

CMD ["./backend"]
