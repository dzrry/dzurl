FROM golang:alpine as builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64
WORKDIR /avito-auto-unit/
COPY go.mod go.sum ./

RUN go mod download
COPY . .

RUN go build -o shortener ./cmd/main.go

FROM scratch

WORKDIR /avito-auto-unit/

COPY --from=builder /avito-auto-unit/shortener shortener
COPY config/ config/

EXPOSE 8080

ENTRYPOINT ["/avito-auto-unit/shortener"]