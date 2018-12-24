# This file creates a container that runs a gofido server
# Important:
# Author: Caliopen
# Date: 2018-12-27

FROM public-registry.caliopen.org/caliopen_go as builder

ADD . /go/src/github.com/CaliOpen/gofido
WORKDIR /go/src/github.com/CaliOpen/gofido

# Fetch dependencies needed for Caliopen GO apps
# RUN govendor sync -v

RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' github.com/CaliOpen/gofido

FROM scratch
MAINTAINER Caliopen

# Add CA certificates
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /go/src/github.com/CaliOpen/gofido/gofido /usr/bin/gofido

WORKDIR /etc/caliopen
ENTRYPOINT ["gofido", "-c" ,"gofido.yaml.template"]

EXPOSE 31415
