docker build \
    --no-cache=true \
    --build-arg BUILD_DATE=$(date -u +'%Y-%m-%dT%H:%M:%SZ') \
    --build-arg VCS_REF=$(git log --pretty=format:'%h' -n 1) \
    --build-arg BUILD_VERSION="1.0.0" \
    -t blaumedia/go-gtm-ga-proxy:latest \
    .

docker push blaumedia/go-gtm-ga-proxy:latest
