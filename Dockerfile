FROM golang:1.17-alpine3.13 as builder
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-w -extldflags "-static"' -o ./url_shortener ./cmd/url_shortener

FROM scratch
WORKDIR /src
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /src/url_shortener .

EXPOSE 9000
CMD ["./url_shortener"]