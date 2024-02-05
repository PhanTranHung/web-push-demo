


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



# Stage 3 - deployment -------------------------------------------
FROM scratch as production


WORKDIR /webpush
COPY . .

COPY --from=builder /webpush/server /webpush/server


COPY --from=builder /lib/aarch64-linux-gnu/libm.so.6    /lib/aarch64-linux-gnu/libm.so.6
COPY --from=builder /lib/aarch64-linux-gnu/libc.so.6    /lib/aarch64-linux-gnu/libc.so.6
COPY --from=builder /lib/ld-linux-aarch64.so.1          /lib/ld-linux-aarch64.so.1

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.cr
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

EXPOSE 8080

CMD ["/webpush/server"]
