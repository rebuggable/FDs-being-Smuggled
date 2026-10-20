include Makefiles/actuateEngine.mk
include Makefiles/createImageOfEngine.mk
include Makefiles/createImageOfEphemeral.mk
include Makefiles/deleteImage.mk
include Makefiles/createExecutable.mk
include Makefiles/dispatch.mk
include Makefiles/actuateRelease.mk
include Makefiles/metadata.mk

.PHONY: build all image run executable clean engine dependencies release version

# make dependencies
dependencies:
	$(GO_BIN) mod download

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

# make release image engine
# make release image ephemeral
# make release image all

# make release executable
# make release all
# make release
release:
	@$(call dispatch,all,actuateExecutableRelease)
	@$(call dispatch,executable,actuateExecutableRelease)

all:
	@:

executable:
	@:

engine:
	@:

version:
	@echo $(VERSION)
