include Makefiles/actuateEngine.mk
include Makefiles/createImageOfEngine.mk
include Makefiles/createImageOfEphemeral.mk
include Makefiles/deleteImage.mk
include Makefiles/createExecutable.mk
include Makefiles/dispatch.mk

.PHONY: build all image run executable clean engine dependencies

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

all:
	@:

executable:
	@:

engine:
	@:
