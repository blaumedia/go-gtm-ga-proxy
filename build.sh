#!/bin/sh
#
# Run:  ./build.sh 1.0.0
# to build image
#
# Run:  ./build.sh 1.0.0 1
# to build & push image
#

if [ $# -eq 0 ]
  then
    echo "Please supply the version number (x.x.x) as first argument."
    echo "Example: ./build.sh 1.0.0"
    exit 1
fi

if [ $3 = "1" ]
  then
    docker build \
    --no-cache=true \
    --build-arg BUILD_DATE=$(date -u +'%Y-%m-%dT%H:%M:%SZ') \
    --build-arg VCS_REF=$(git log --pretty=format:'%h' -n 1) \
    --build-arg BUILD_VERSION="$1" \
    -t blaumedia/go-gtm-ga-proxy:latest-pluginengine \
    -f Dockerfile-PluginEngine.Dockerfile \
    .
else
  docker build \
    --no-cache=true \
    --build-arg BUILD_DATE=$(date -u +'%Y-%m-%dT%H:%M:%SZ') \
    --build-arg VCS_REF=$(git log --pretty=format:'%h' -n 1) \
    --build-arg BUILD_VERSION="$1" \
    -t blaumedia/go-gtm-ga-proxy:latest \
    -f Dockerfile \
    .
fi

if [ $2 = "1" ]
  then
    docker push blaumedia/go-gtm-ga-proxy:latest
    docker tag blaumedia/go-gtm-ga-proxy:latest docker.pkg.github.com/blaumedia/go-gtm-ga-proxy/go-gtm-ga-proxy:latest
    docker push docker.pkg.github.com/blaumedia/go-gtm-ga-proxy/go-gtm-ga-proxy:latest

    docker tag blaumedia/go-gtm-ga-proxy:latest blaumedia/go-gtm-ga-proxy:$1
    docker push blaumedia/go-gtm-ga-proxy:$1

    docker tag blaumedia/go-gtm-ga-proxy:latest docker.pkg.github.com/blaumedia/go-gtm-ga-proxy/go-gtm-ga-proxy:$1
    docker push docker.pkg.github.com/blaumedia/go-gtm-ga-proxy/go-gtm-ga-proxy:$1
else
  echo "Skipping push to DockerHub."
fi