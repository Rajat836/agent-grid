# Args
# GO_VERSION: Go version to use
# UID: User ID for the appuser
# PORT: Port to expose

ARG GO_VERSION=1.24-bookworm
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION} AS base

RUN apt-get install ca-certificates -y

ARG UID=10001
RUN adduser \
  --disabled-password \
  --gecos "" \
  --home "/nonexistent" \
  --shell "/sbin/nologin" \
  --no-create-home \
  --uid "$UID" \
  appuser

WORKDIR $GOPATH/app/
COPY . .

# Configuring git to use github credentials
ENV GOPRIVATE=bitbucket.org/risk-manager-x
ARG BITBUCKET_USER
ARG BITBUCKET_PASSWORD
RUN echo "machine bitbucket.org login ${BITBUCKET_USER} password ${BITBUCKET_PASSWORD}" > ~/.netrc
RUN chmod 600 ~/.netrc

ARG COMMIT_HASH
ENV COMMIT_HASH=$COMMIT_HASH

# Download dependencies
RUN go mod download
RUN go mod verify

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GODEBUG=gctrace=1 \
    go build -p 2 -ldflags="-w -s -X main.Version=$COMMIT_HASH" \
    -trimpath -o /bin/server ./cmd/server


FROM scratch

COPY --from=base /bin/server /app/server
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=base /go/app/migrations/postgres ./migrations/postgres

# Expose the port that the application listens on.
ARG PORT=4100
EXPOSE "$PORT"

ENTRYPOINT [ "./app/server", "--config", "/app/config/config.yml" ]
