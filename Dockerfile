FROM golang:1.25.1-trixie@sha256:61226c61f37cb86253c4ac486ef22c47f14bfddb8f60bb4805bfc165001be758

WORKDIR /hambot

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build --tags fts5

CMD ["./hambot"]
