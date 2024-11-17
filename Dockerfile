FROM golang:1.23.3

WORKDIR /app

COPY testik/go.mod testik/go.sum ./
RUN go mod download && go mod verify

COPY testik/ ./

RUN go build -v -o ./dns-proxy .

CMD ["/app/dns-proxy"]