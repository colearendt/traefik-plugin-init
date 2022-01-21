# from https://stackoverflow.com/a/46532352/6570011
# build stage
FROM golang:1.17
LABEL org.opencontainers.image.source=https://github.com/colearendt/traefik-plugin-init

ADD . /src
RUN set -x && \
    cd /src && \
    CGO_ENABLED=0 GOOS=linux go build -a -o traefik-plugin-init

# final stage
FROM alpine
WORKDIR /app
COPY --from=0 /src/traefik-plugin-init /app/
ENTRYPOINT /app/traefik-plugin-init
