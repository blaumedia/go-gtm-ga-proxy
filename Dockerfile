#
# Building the Go library using multistage builds
#
FROM golang:1.14.2
WORKDIR /go/src/GoGtmGaProxy/
COPY ./server/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

#
# Copying the built binary to alpine image
#
FROM alpine:latest

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

# Using port 8080 because we won't run the application as root
EXPOSE 8080

RUN apk --no-cache add ca-certificates

RUN addgroup -S docker -g 433 && \
    adduser -u 431 -S -g docker -h /app -s /sbin/nologin docker && \
    chown -R docker:docker /app

USER docker
WORKDIR /app/
COPY --from=0 /go/src/GoGtmGaProxy/app GoGtmGaProxy
CMD ["./GoGtmGaProxy"]