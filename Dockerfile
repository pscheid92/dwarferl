FROM golang:1.19beta1-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app
COPY . .

RUN go mod download
RUN go mod verify

RUN CGO_ENABLED=0 go build -tags=go_json -ldflags='-w -s'

FROM scratch

COPY ./assets /assets
COPY ./templates /templates

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/dwarferl /dwarferl

EXPOSE 8080
ENTRYPOINT ["/dwarferl"]
