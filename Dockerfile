#
# ===========
# Build image
# ===========
#
FROM golang:1.19 as builder
COPY . /go/src/github.com/registrobr/rdap-client
WORKDIR /go/src/github.com/registrobr/rdap-client
RUN mkdir /apps
RUN go build -mod=vendor -ldflags="-w -s" -o /apps/rdap-client

#
# ====================
# Final delivery image
# ====================
#
FROM debian:stable-slim

COPY --from=builder /apps/* /apps/

RUN apt update \
    && apt -y install ca-certificates \
    && rm -rf /var/lib/apt/lists/*

ARG BUILD_DATE
ARG BUILD_VCS_REF
ARG BUILD_VERSION

ENV API_VERSION ${BUILD_VCS_REF}

LABEL org.label-schema.build-date=$BUILD_DATE \
      org.label-schema.description="RDAP client" \
      org.label-schema.name="rdap-client" \
      org.label-schema.schema-version="1.0" \
      org.label-schema.url="https://registro.br" \
      org.label-schema.vcs-url="https://github.com/registrobr/rdap-client" \
      org.label-schema.vcs-ref=$BUILD_VCS_REF \
      org.label-schema.vendor="NIC.br" \
      org.label-schema.version=$BUILD_VERSION

ENTRYPOINT ["/apps/rdap-client"]
