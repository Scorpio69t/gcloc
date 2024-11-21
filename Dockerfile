FROM golang:1.20-buster AS builder

RUN apt-get update \
 && apt-get install -y --no-install-recommends \
  upx-ucl

WORKDIR /build

COPY . .

ARG VERSION=dev
ARG GIT_COMMIT=none
ARG BUILD_DATE=unknown

RUN GO111MODULE=on CGO_ENABLED=0 go build \
      -ldflags='-w -s -extldflags "-static" -X github.com/Scorpio69t/gcloc/cmd.Version=${VERSION} -X github.com/Scorpio69t/gcloc/cmd.GitCommit=${GIT_COMMIT} -X github.com/Scorpio69t/gcloc/cmd.BuildDate=${BUILD_DATE}' \
      -o ./bin/gcloc app/gcloc/main.go \
 && upx-ucl --best --ultra-brute ./bin/gcloc

FROM scratch

COPY --from=builder /build/bin/gcloc /bin/

WORKDIR /workdir

ENTRYPOINT ["/bin/gcloc"]
