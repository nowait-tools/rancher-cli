# build
FROM golang:1.8-alpine AS build_env
WORKDIR /go/src/github.com/nowait-tools/rancher-cli
ADD  . .
RUN  go build -x -o rancher-cli .

# runtime
FROM alpine
RUN apk --update add ca-certificates
COPY --from=build_env /go/src/github.com/nowait-tools/rancher-cli /
RUN chmod +x rancher-cli
ENTRYPOINT ["/rancher-cli"]
