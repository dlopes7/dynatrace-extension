FROM golang:1.16 as builder

WORKDIR /workspace

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY cmd/ cmd/
COPY pkg/ pkg/

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o extension-watcher cmd/extension-watcher/watcher.go

FROM gcr.io/distroless/static:nonroot

WORKDIR /
COPY --from=builder /workspace/extension-watcher .
USER 0:0
ENTRYPOINT ["/extension-watcher"]
