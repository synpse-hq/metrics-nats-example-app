FROM quay.io/synpse/golang:1.16
ARG version
WORKDIR /app

# <- COPY go.mod and go.sum files to the workspace
COPY go.mod . 
COPY go.sum .

# Get dependancies - will also be cached if we won't change mod/sum
RUN go mod download

# COPY the source code as the last step
COPY . .

ARG GOARCH=amd64

RUN make app

FROM quay.io/synpse/alpine:3.9
RUN apk --update add git openssh tar gzip ca-certificates \
  bash curl
ARG version
COPY --from=0 /app/release/app/app /bin/
ENTRYPOINT ["/bin/app"]
