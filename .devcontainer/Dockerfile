FROM golang:1.12 as dev
WORKDIR /app
ENV SHELL /bin/bash

RUN go get -v -u \
    golang.org/x/tools/cmd/gopls \
    github.com/onsi/ginkgo/ginkgo  \
    github.com/onsi/gomega/...  \
    github.com/ramya-rao-a/go-outline \
    github.com/mdempsky/gocode \
    github.com/uudashr/gopkgs/cmd/gopkgs \
    github.com/ramya-rao-a/go-outline \
    github.com/acroca/go-symbols \
    golang.org/x/tools/cmd/guru \
    golang.org/x/tools/cmd/gorename \
    github.com/go-delve/delve/cmd/dlv \
    github.com/stamblerre/gocode \
    github.com/rogpeppe/godef \
    golang.org/x/tools/cmd/goimports \
    golang.org/x/lint/golint \
    golang.org/x/tools/cmd/gopls \
    github.com/cweill/gotests/... \
    github.com/fatih/gomodifytags \
    github.com/josharian/impl \
    github.com/davidrjenni/reftools/cmd/fillstruct \
    github.com/haya14busa/goplay/cmd/goplay \
    github.com/godoctor/godoctor 

run chown -R 1000 /go

