FROM scratch

COPY ./assets /assets
COPY ./templates /templates
COPY ./dwarferl /dwarferl

EXPOSE 8080
ENTRYPOINT ["/dwarferl"]
