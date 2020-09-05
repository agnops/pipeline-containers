FROM alpine:3.8 as build

RUN apk add --update --no-cache ca-certificates git

ARG HELM_VERSION
ARG HELM_SHA256SUM
ARG KUBECTL_VERSION

ENV FILENAME=helm-${HELM_VERSION}-linux-amd64.tar.gz

ADD https://storage.googleapis.com/kubernetes-release/release/v1.16.5/bin/linux/amd64/kubectl /usr/local/bin/kubectl
RUN apk add --update -t deps curl tar gzip
RUN curl -L https://get.helm.sh/${FILENAME} > ${FILENAME} && \
    echo "${HELM_SHA256SUM}  ${FILENAME}" > helm_${HELM_VERSION}_SHA256SUMS && \
    sha256sum -cs helm_${HELM_VERSION}_SHA256SUMS && \
    tar zxv -C /tmp -f ${FILENAME} && \
    rm -f ${FILENAME}

FROM alpine:latest

RUN apk add --update coreutils

COPY --from=build /tmp/linux-amd64/helm /bin/helm
COPY --from=build /usr/local/bin/kubectl /usr/local/bin/kubectl
RUN chmod +x /usr/local/bin/kubectl