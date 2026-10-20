CURRENT_COMMIT := $(shell git rev-parse --short --verify HEAD 2>/dev/null)
HAS_COMMITED := $(if $(CURRENT_COMMIT),True,False)

# Have we commited a change?
ifeq ($(HAS_COMMITED),True)
  CURRENT_TAG := $(shell git describe --exact-match --tags $(CURRENT_COMMIT) 2>/dev/null)
  PREVIOUS_TAG := $(shell git describe --tags --abbrev=0 2>/dev/null)

  IS_CURRENTLY_TAGGED := $(if $(CURRENT_TAG),True,False)
  BEEN_PREVIOUSLY_TAGGED := $(if $(PREVIOUS_TAG),True,False)

  ifeq ($(IS_CURRENTLY_TAGGED),True)
    RAW_VERSION := $(CURRENT_TAG)
  else ifeq ($(BEEN_PREVIOUSLY_TAGGED),True)
		DISTANCE := $(shell git rev-list $(PREVIOUS_TAG)..$(CURRENT_COMMIT) --count 2>/dev/null)
    RAW_VERSION := $(PREVIOUS_TAG)-$(DISTANCE).$(CURRENT_COMMIT)
  else
    RAW_VERSION := v0.0.0-$(CURRENT_COMMIT)
  endif

# Otherwise (didn't commit a change yet)
else
  IS_CURRENTLY_TAGGED := False
  BEEN_PREVIOUSLY_TAGGED := False
  RAW_VERSION := v0.0.0
endif

VERSION = $(patsubst v%,%,$(RAW_VERSION))
