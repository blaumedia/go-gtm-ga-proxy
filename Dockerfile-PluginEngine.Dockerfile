#
# Building the Go binary using multistage builds
#
FROM golang:1.14.4-alpine3.12
RUN apk --no-cache add gcc g++

COPY server /go/src/server
RUN GO111MODULE=on go get -v github.com/tdewolff/minify/v2 && \
  cd /go/src/server && GO111MODULE=on CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o app .

COPY ./plugins/*.go /go/src/server/plugins/
RUN cd /go/src/server/plugins \
    && CGO_ENABLED=1 GOOS=linux find . -name \*.go -type f -exec go build -a -buildmode=plugin -o $(basename {} .go).so {} \;

#
# Copying the built binary to alpine image
#
FROM golang:1.14.4-alpine3.12

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

RUN apk --no-cache add ca-certificates uglify-js

# RUN useradd -u 431 -r -d /app -s /sbin/nologin docker
RUN addgroup -S docker -g 433 && \
    adduser -u 431 -S -g docker -h /app -s /sbin/nologin docker

USER docker
WORKDIR /app/

COPY --chown=docker:docker --from=0 /go/src/server/app GoGtmGaProxy
COPY --chown=docker:docker --from=0 /go/src/server/plugins/*.so ./plugins/
CMD ["./GoGtmGaProxy"]