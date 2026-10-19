eq = $(if $(subst $(1),,$(2))$(subst $(2),,$(1)),,True)

arguments = $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))

define dispatch
	$(if $(call eq,$(strip $(1)),$(strip $(arguments))),$(call $(2)),:)
endef
