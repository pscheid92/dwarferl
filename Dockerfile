FROM alpine:latest

RUN apk add --no-cache libc6-compat

COPY ./dwarferl /dwarferl

EXPOSE 8080
ENTRYPOINT ["/dwarferl"]
