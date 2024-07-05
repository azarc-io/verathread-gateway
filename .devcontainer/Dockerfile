# this image is built and updated weekly
# https://github.com/cse-labs/codespaces-images

# FROM mcr.microsoft.com/devcontainers/go:1.22-bookworm

FROM golang:1.22.4-bullseye

RUN go install github.com/fatih/gomodifytags@latest && \
    go install github.com/cweill/gotests/gotests@latest && \
    go install github.com/go-delve/delve/cmd/dlv@latest && \
    go install golang.org/x/tools/gopls@latest && \
    go install honnef.co/go/tools/cmd/staticcheck@latest && \
    go install github.com/josharian/impl@latest && \
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest