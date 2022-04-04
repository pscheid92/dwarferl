FROM alpine:latest

RUN apk add --no-cache libc6-compat

COPY ./dwarferl /dwarferl
COPY ./templates /templates

EXPOSE 8080
ENTRYPOINT ["/dwarferl"]
