FROM --platform=${BUILDPLATFORM} quay.io/synpse/alpine:3.9
RUN apk --update add git openssh tar gzip ca-certificates \
  bash curl
ARG TARGETPLATFORM
ARG BUILDPLATFORM

COPY ./release/${BUILDPLATFORM}/app /bin/
ENTRYPOINT ["/bin/app"]
