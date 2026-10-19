include Makefiles/binaries.mk

define createImageOfEngine
	$(DOCKER_BIN) build -f ./Dockerfiles/engine/engine.Dockerfile -t a-mule-engine .
endef
