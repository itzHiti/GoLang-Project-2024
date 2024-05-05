FROM golang:1.21.6-alpine

WORKDIR /visual/ocm

COPY . .

RUN go mod download

RUN go build -o OCM ./cmd/OCM

EXPOSE 8081

CMD ["./ocm"]
