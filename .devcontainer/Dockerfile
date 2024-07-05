# this image is built and updated weekly
# https://github.com/cse-labs/codespaces-images

FROM mcr.microsoft.com/devcontainers/base

RUN apt-get update
RUN apt-get install -y apt-transport-https ca-certificates curl gnupg
RUN curl -fsSL https://pkgs.k8s.io/core:/stable:/v1.30/deb/Release.key | gpg --dearmor -o /etc/apt/keyrings/kubernetes-apt-keyring.gpg
RUN chmod 644 /etc/apt/keyrings/kubernetes-apt-keyring.gpg
RUN echo 'deb [signed-by=/etc/apt/keyrings/kubernetes-apt-keyring.gpg] https://pkgs.k8s.io/core:/stable:/v1.30/deb/ /' | tee /etc/apt/sources.list.d/kubernetes.list
RUN chmod 644 /etc/apt/sources.list.d/kubernetes.list
RUN apt-get update
RUN apt-get install -y kubectl

#FROM golang:1.22.4-bullseye

#RUN go install github.com/fatih/gomodifytags@latest && \
#    go install github.com/cweill/gotests/gotests@latest && \
#    go install github.com/go-delve/delve/cmd/dlv@latest && \
#    go install golang.org/x/tools/gopls@latest && \
#    go install honnef.co/go/tools/cmd/staticcheck@latest && \
#    go install github.com/josharian/impl@latest && \
#    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
