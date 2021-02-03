FROM golang:1.14-buster AS build
ENV GOPROXY=https://proxy.golang.org
WORKDIR /go/src/github.com/brito-rafa/kubectl-mutate-dc2deployment
COPY . .
RUN CGO_ENABLED=0 go build -o /go/bin/kubectl-mutate-dc2deployment .


FROM ubuntu:bionic
USER nobody:nogroup
#ENTRYPOINT ["/bin/bash", "-c", "cp /plugins/* /target/."]
