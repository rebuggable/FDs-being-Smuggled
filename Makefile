GO_BIN := $(shell command -v go)
DOCKER_BIN := $(shell command -v docker)

GOOS ?= linux
GOARCH ?= amd64

.PHONY: build all image run executable clean engine

eq = $(if $(subst $(1),,$(2))$(subst $(2),,$(1)),,True)

arguments = $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))

define dispatch
	$(if $(call eq,$(strip $(1)),$(strip $(arguments))),$(call $(2)),:)
endef

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

define createImageOfEphemeral
	$(DOCKER_BIN) build -f ./Dockerfiles/ephemeral/ephemeral.Dockerfile -t a-mule-ephemeral .
endef

define createImageOfEngine
	$(DOCKER_BIN) build -f ./Dockerfiles/engine/engine.Dockerfile -t a-mule-engine .
endef

define actuateEngine
	$(DOCKER_BIN) network inspect mule-network >/dev/null 2>&1 || $(DOCKER_BIN) network create mule-network ; \
	$(DOCKER_BIN) run -d \
		--privileged \
		--name mule-engine \
		--network mule-network \
		--network-alias docker \
		-e DOCKER_TLS_CERTDIR=/certs \
		-v mule-docker-certs-ca:/certs/ca \
		-v mule-docker-certs-client:/certs/client \
		a-mule-engine
endef

define deleteImage
	$(DOCKER_BIN) image inspect $(1) >/dev/null 2>&1 && $(DOCKER_BIN) image rm $(1)
endef

# make build all
# make build executable
build:
	@$(call dispatch,all,createExecutable)
	@$(call dispatch,all,createImageOfEphemeral)
	@$(call dispatch,all,createImageOfEngine)

	@$(call dispatch,executable,createExecutable)

# make image all
# make image ephemeral
# make image engine
image:
	@$(call dispatch,all,createImageOfEphemeral)
	@$(call dispatch,all,createImageOfEngine)

	@$(call dispatch,ephemeral,createImageOfEphemeral)
	@$(call dispatch,engine,createImageOfEngine)

# make run engine
run:
	@$(call dispatch,engine,actuateEngine)

# make clean
clean:
	rm -r ./dest ; \
	$(call deleteImage, a-mule-ephemeral) ; \
	$(call deleteImage, a-mule-engine)

all:
	@:

executable:
	@:

engine:
	@: