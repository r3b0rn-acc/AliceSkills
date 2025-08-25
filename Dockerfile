FROM --platform=$BUILDPLATFORM golang:1.24-alpine AS build

RUN apk add --no-cache ca-certificates

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG PKG=./cmd/skillsrv

ARG TARGETOS
ARG TARGETARCH

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go build $GOFLAGS -trimpath -buildvcs=false -ldflags="-s -w" \
      -o /out/skillsrv ${PKG}

FROM scratch AS final

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=build /out/skillsrv /skillsrv

ENTRYPOINT ["/skillsrv"]
