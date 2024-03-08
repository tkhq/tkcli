FROM stagex/busybox:latest@sha256:a925ae77d4c8e633bc1706e23543a4b1b3a2b183257ed840b2c77aaaadeb3465 AS busybox
FROM stagex/ca-certificates:latest@sha256:c6f95ed1ae7ec9e6e19d24fc4569bf46ae886042ebc4a16cc8ec53b762290379 AS ca-certificates
FROM stagex/go:v1.21.4@sha256:26c3aab06c302fe6397793cc18608993834cdaadd96c059d214a96929af9baaf AS go

FROM scratch AS base
LABEL org.opencontainers.image.source https://github.com/tkhq/tkcli
COPY --from=ca-certificates . /
ENTRYPOINT ["/app"]

FROM base AS fetch
ENV GOPROXY="https://proxy.golang.org,direct"
COPY --from=busybox . /
COPY --from=go . /

FROM fetch AS build
RUN go env -w GOMODCACHE=/go/pkg/mod
RUN go env -w GOSUMDB=sum.golang.org
ARG BUILD_SRC_DIR
ARG GO_BUILDFLAGS="-trimpath -buildvcs=false"
ARG GO_LDFLAGS="-s -w -buildid= -extldflags=-static"

COPY . /home/user/tkcli

# Prefetch go modules
FROM build AS fetch_modules
WORKDIR /home/user/tkcli
RUN --mount=type=cache,target=/go/pkg/mod go mod download -json

# Turnkey
FROM build AS turnkey_build
WORKDIR /home/user/tkcli/src/cmd/turnkey
RUN go env GOMODCACHE
RUN --mount=type=cache,target=/go/pkg/mod go build \
  ${GO_BUILDFLAGS} \
  -ldflags="${GO_LDFLAGS}" \
  -o /home/user/bin/app .

FROM base AS turnkey
ARG LABEL
LABEL org.opencontainers.image.title ${LABEL}
COPY --from=turnkey_build /home/user/bin/app /app
USER 100:100