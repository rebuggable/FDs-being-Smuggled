include Makefiles/binaries.mk
include Makefiles/dispatch.mk

GOOS ?= $(shell $(GO_BIN) env GOHOSTOS)
GOARCH ?= $(shell $(GO_BIN) env GOHOSTARCH)

define createExecutable
	$(let _GOOS _GOARCH, $(or $(1),$(GOOS)) $(or $(2),$(GOARCH)), \
	CGO_ENABLED=0 \
	GOOS=$(_GOOS) \
	GOARCH=$(_GOARCH) \
	$(GO_BIN) build \
		-trimpath \
		-ldflags="-s -w" \
		-o ./dest/mule-$(_GOOS)-$(_GOARCH)$(if $(call eq, $(_GOOS), windows),.exe) \
		./interlocutor)
endef
