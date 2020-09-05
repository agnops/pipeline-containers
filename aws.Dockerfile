FROM alpine:3.8 as build

RUN apk add --update --no-cache ca-certificates git

ARG HELM_VERSION
ARG HELM_SHA256SUM
ARG KUBECTL_VERSION

ENV FILENAME=helm-${HELM_VERSION}-linux-amd64.tar.gz

ADD https://storage.googleapis.com/kubernetes-release/release/${KUBECTL_VERSION}/bin/linux/amd64/kubectl /usr/local/bin/kubectl
RUN apk add --update -t deps curl tar gzip
RUN curl -L https://get.helm.sh/${FILENAME} > ${FILENAME} && \
    echo "${HELM_SHA256SUM}  ${FILENAME}" > helm_${HELM_VERSION}_SHA256SUMS && \
    sha256sum -cs helm_${HELM_VERSION}_SHA256SUMS && \
    tar zxv -C /tmp -f ${FILENAME} && \
    rm -f ${FILENAME}

FROM alpine:latest

COPY --from=build /tmp/linux-amd64/helm /bin/helm
COPY --from=build /usr/local/bin/kubectl /usr/local/bin/kubectl

ADD https://github.com/kubernetes-sigs/aws-iam-authenticator/releases/download/v0.5.0/aws-iam-authenticator_0.5.0_linux_amd64 /usr/local/bin/aws-iam-authenticator

RUN set -x
RUN apk add --update --no-cache curl ca-certificates python3 py-pip
RUN chmod +x /usr/local/bin/kubectl && chmod +x /usr/local/bin/aws-iam-authenticator
RUN pip install --upgrade awscli

RUN aws --version
RUN kubectl version --client
RUN helm version