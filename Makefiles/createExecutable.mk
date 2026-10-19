include Makefiles/binaries.mk

GOOS ?= linux
GOARCH ?= amd64

define createExecutable
	CGO_ENABLED=0 \
	GOOS=$(GOOS) \
	GOARCH=$(GOARCH) \
	$(GO_BIN) build \
		-trimpath \
		-ldflags="-s -w" \
		-o ./dest/mule \
		./interlocutor
endef
