FROM golang:1.19 AS builder
WORKDIR /go/src/open-cluster-management.io/addon-contrib/example-addon
COPY . .
ENV GO_PACKAGE open-cluster-management.io/addon-contrib/example-addon

ENV GOFLAGS ""
RUN go env
RUN go build -a -o example-addon main.go

FROM registry.access.redhat.com/ubi8/ubi-minimal:latest
COPY --from=builder /go/src/open-cluster-management.io/addon-contrib/example-addon /

RUN microdnf update && microdnf clean all

ENTRYPOINT ["/example-addon"]
