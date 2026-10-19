include Makefiles/binaries.mk

define createImageOfEphemeral
	$(DOCKER_BIN) build -f ./Dockerfiles/ephemeral/ephemeral.Dockerfile -t a-mule-ephemeral .
endef
