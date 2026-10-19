include Makefiles/createExecutable.mk

define actuateExecutableRelease
	$(call createExecutable, darwin, amd64)
	$(call createExecutable, darwin, arm64)
	$(call createExecutable, linux, amd64)
	$(call createExecutable, linux, arm)
	$(call createExecutable, linux, arm64)
	$(call createExecutable, linux, ppc64le)
	$(call createExecutable, linux, riscv64)
	$(call createExecutable, linux, s390x)
	$(call createExecutable, windows, amd64)
	$(call createExecutable, windows, arm64)
endef
