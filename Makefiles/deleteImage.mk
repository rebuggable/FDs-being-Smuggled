include Makefiles/binaries.mk

define deleteImage
	$(DOCKER_BIN) image inspect $(1) >/dev/null 2>&1 && $(DOCKER_BIN) image rm $(1)
endef
