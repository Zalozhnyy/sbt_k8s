# Stage 1 - Binary Build
# BUILD_X args should be passed at build time as docker build args
FROM golang:1.22-alpine3.10 AS builder
ARG BIN_VERSION
ARG GIT_COMMIT
ARG GIT_SHA
ARG GIT_TAG
ARG GIT_DIRTY
ENV BIN_OUTDIR=./
ENV BIN_NAME=sbt_k8s
RUN apk update && apk add build-base git libressl-dev
WORKDIR /usr/src/sbt_k8s
# install dependencies in separate docker layer
COPY go.mod .
COPY go.sum .
RUN go mod download
# copy application source and build
COPY ./ .
RUN make static

# Stage 2 - Final Image
# The application should be statically linked
FROM alpine:3.10
RUN apk update \
	&& apk add --no-cache ca-certificates \
	&& rm -rf /var/cache/apk/* \
    && addgroup sbt_k8s \
	&& adduser -D -H -G sbt_k8s sbt_k8s
COPY --from=builder /usr/src/sbt_k8s/sbt_k8s /usr/bin/sbt_k8s
VOLUME ["/etc/sbt_k8s"]
ENTRYPOINT ["sbt_k8s"]
EXPOSE 8080
USER sbt_k8s
