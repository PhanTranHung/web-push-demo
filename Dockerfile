


# Stage 1 - building -------------------------------------------
FROM golang:1.20 AS builder

WORKDIR /webpush
COPY ./go.mod /webpush
COPY ./go.sum /webpush
RUN --mount=type=cache,target=/var/cache/apt go mod download


COPY . .

RUN --mount=type=cache,target=/var/cache/apt VER=0.0.0-beta && \
    CURRENT_DATE=$(date -u "+%Y-%m-%dT%H:%M:%S") && \
    CGO_ENABLED=1 GOOS=linux        \
    go build -o server -ldflags="-X trada/config.bVersion=${VER} -X trada/config.bDateTime=${CURRENT_DATE}" \
    ./main/main.go

# Create admin user
ENV USER=admin
ENV UID=10001
RUN adduser --disabled-password --gecos "" --no-create-home --shell "/sbin/nologin" --uid "${UID}" "${USER}"

CMD ["/webpush/server"]
