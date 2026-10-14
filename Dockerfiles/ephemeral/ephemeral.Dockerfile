FROM golang:1.27.1-alpine AS Build

WORKDIR /opt/FDs-being-Smuggled

COPY go.mod .
COPY go.sum .

COPY . .

RUN CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    go build \
      -trimpath \
      -ldflags="-s -w" \
      -o /usr/local/bin/orchestrator \
      ./orchestrator

FROM scratch AS Final

COPY --from=Build /usr/local/bin/orchestrator /usr/local/bin/orchestrator

ENV DOCKER_MODE=DooD

ENTRYPOINT ["/usr/local/bin/orchestrator"]
