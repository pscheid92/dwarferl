FROM golang:1.18-alpine AS builder

RUN apk add --no-cache git libc6-compat

WORKDIR /app
COPY . .

RUN go mod download
RUN go mod verify

RUN CGO_ENABLED=0 go build -ldflags='-w -s'

FROM scratch

COPY ./assets /assets
COPY ./templates /templates

COPY --from=builder /app/dwarferl /dwarferl

EXPOSE 8080
ENTRYPOINT ["/dwarferl"]
