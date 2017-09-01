FROM golang:1.8

COPY . /go/src/github.com/anchorfree/prometheus-ha-proxy
RUN curl https://glide.sh/get | sh \
    && cd /go/src/github.com/anchorfree/prometheus-ha-proxy \
    && glide install -v
RUN cd /go/src/github.com/anchorfree/prometheus-ha-proxy \
    && CGO_ENABLED=0 go build -o /build/prometheus-ha-proxy  *.go


FROM alpine

RUN apk add curl --update-cache
COPY --from=0 /build/prometheus-ha-proxy /prometheus-ha-proxy

EXPOSE 9374

ENTRYPOINT ["/prometheus-ha-proxy"]
