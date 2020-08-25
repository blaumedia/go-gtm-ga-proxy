#
# Building the Go binary using multistage builds
#
FROM golang:1.14.4-alpine3.12
COPY server /go/src/server
RUN GO111MODULE=on go get -v github.com/tdewolff/minify/v2 && \
  cd /go/src/server && GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

#
# Copying the built binary to alpine image
#
FROM alpine:3.12

ARG BUILD_DATE
ARG BUILD_VERSION
ARG VCS_REF

LABEL author="Dennis Paul"
LABEL maintainer="dennis@blaumedia.com"

LABEL org.label-schema.schema-version="1.0"
LABEL org.label-schema.build_date=$BUILD_DATE
LABEL org.label-schema.name="Go-GTM-GA-Proxy"
LABEL org.label-schema.vcs-url="https://github.com/blaumedia/go-gtm-ga-proxy"
LABEL org.label-schema.version=$BUILD_VERSION
LABEL org.label-schema.vcs-ref=$VCS_REF

ENV APP_VERSION ${BUILD_VERSION}

# Using port 8080 because we won't run the application as root
EXPOSE 8080

RUN apk --no-cache add ca-certificates uglify-js curl

RUN addgroup -S docker -g 433 && \
    adduser -u 431 -S -g docker -h /app -s /sbin/nologin docker

USER docker

WORKDIR /app/
COPY --from=0 /go/src/server/app GoGtmGaProxy

HEALTHCHECK --interval=10s --timeout=3s --start-period=5s --retries=3 CMD curl -f http://localhost:8080/$JS_SUBDIRECTORY/$GA_FILENAME || exit 1

CMD ["./GoGtmGaProxy"]