FROM golang:1.27.3-trixie@sha256:c124b9ad3452f46061b5054607d26ea51730dae5372c5e469920063ba2bc6dcc

WORKDIR /hambot

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build --tags fts5

CMD ["./hambot"]
