# We need wget for a health check
FROM busybox AS busybox

ARG BUSYBOX_VERSION=1.31.0-i686-uclibc
ADD https://busybox.net/downloads/binaries/$BUSYBOX_VERSION/busybox_WGET /wget
RUN chmod a+x /wget

# Start by building the application.
FROM golang:1.15 as build

WORKDIR /go/src/app
ADD . /go/src/app

RUN go get -d -v ./...

RUN CGO_ENABLED=0 go build -o /go/bin/app -ldflags '-extldflags "-static"'

# Now copy it into our base image.
FROM scratch
LABEL maintainer="Anze Jensterle <anze@anze.dev>"
LABEL Description="This image is used to start the ToDo API"
LABEL org.opencontainers.image.source=https://github.com/craftbyte/msi-k8s
USER 999
COPY --from=build /go/bin/app /
COPY --from=busybox /wget /wget
ENV GIN_MODE=release
ENV PORT=8080
EXPOSE 8080
CMD ["/app"]