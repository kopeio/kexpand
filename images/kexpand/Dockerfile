FROM alpine:3.5

# KUBECTL_SOURCE: Change to kubernetes-dev/ci for CI
ARG KUBECTL_SOURCE=kubernetes-release/release

# KUBECTL_TRACK: Currently latest from KUBECTL_SOURCE. Change to latest-1.3.txt, etc. if desired.
ARG KUBECTL_TRACK=stable.txt

RUN set -ex && \
        apk add --no-cache --virtual deps curl jq ca-certificates && \
        KUBECTL_VERSION=$(curl -SsL --retry 5 "https://storage.googleapis.com/${KUBECTL_SOURCE}/${KUBECTL_TRACK}") && \
        curl -SsL --retry 5 "https://storage.googleapis.com/${KUBECTL_SOURCE}/${KUBECTL_VERSION}/bin/linux/amd64/kubectl" > /usr/bin/kubectl && \
        chmod +x /usr/bin/kubectl && \
        apk del deps && \
        rm -rf /var/cache/apk/*

COPY /.build/artifacts/kexpand /kexpand

CMD /kexpand

