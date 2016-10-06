FROM golang:1.7-alpine
MAINTAINER Nowait <devops@nowait.com>

WORKDIR /src
COPY rancher-cli rancher-cli
RUN chmod +x rancher-cli

ENTRYPOINT ["/src/rancher-cli"]
