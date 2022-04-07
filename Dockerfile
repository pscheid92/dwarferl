FROM alpine:latest

RUN apk add --no-cache libc6-compat

COPY ./assets /assets
COPY ./templates /templates
COPY ./dwarferl /dwarferl

EXPOSE 8080
ENTRYPOINT ["/dwarferl"]
