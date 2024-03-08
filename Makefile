VERSION ?= $(shell git describe --tag --always --dirty --match 'v[0-9]*' --abbrev=8)

KEYS := \
	6B61ECD76088748C70590D55E90A401336C8AAA9 \
	A8864A8303994E3A18ACD1760CAB4418C834B102 \
	66039AA59D823C8BD68DB062D3EC673DF9843E7B \
	DE050A451E6FAF94C677B58B9361DEC647A087BD

LOCAL_BUILD_DIR := 'build'
SRC_DIR := src
KEY_DIR := fetch/keys
OUT_DIR := out

normarch = $(subst arm64,aarch64,$(subst amd64,x86_64,$1))
HOST_ARCH := $(call normarch,$(call lc,$(shell uname -m)))
HOST_ARCH_ALT := $(call altarch,$(HOST_ARCH))
HOST_OS := $(call lc,$(shell uname -s))

GIT_REF := $(shell git log -1 --format=%H)
GIT_AUTHOR := $(shell git log -1 --format=%an)
GIT_KEY := $(shell git log -1 --format=%GP)
GIT_TIMESTAMP := $(shell git log -1 --format=%cd --date=iso)

REGISTRY := local
BUILDER := $(shell which docker)

# Build package with chosen $(BUILDER)
# Supported BUILDERs: docker
# Usage: $(call build,core/$(NAME),$(VERSION),$(TARGET),$(EXTRA_ARGS))
# Notes:
# - Packages are expected to use the following layer names in order:
#   - "fetch": [optional] obtain any artifacts from the internet.
#   - "build": [optional] do any required build work
#   - "package": [required] scratch layer exporting artifacts for distribution
#   - "test": [optional] define any tests
# - Packages may prefix layer names with "text-" if more than one is desired
# - VERSION will be set as a build-arg if defined, otherwise it is "latest"
# - TARGET defaults to "package"
# - EXTRA_ARGS will be blindly injected
# - packages may also define a "test" layer
# - the ulimit line is to workaround a bug in patch when the nofile limit is too large:
#      https://savannah.gnu.org/bugs/index.php?62958
#  TODO:
# - try to disable networking on fetch layers with something like:
#   $(if $(filter fetch,$(lastword $(subst -, ,$(TARGET)))),,--network=none)
# - actually output OCI files for each build (vs plain tar)
# - output manifest.txt of all tar/digest hashes for an easy git diff
# - support buildah and podman
define build
	$(eval $(call determine_platform,$*))
	echo PLATFORM is $(PLATFORM)
	$(eval LANGUAGE := go)
	$(eval NAME := $(2))
	$(eval VERSION := $(if $(3),$(3),latest))
	$(eval TARGET := $(if $(4),$(4),package))
	$(eval EXTRA_ARGS := $(if $(5),$(5),))
	$(eval BUILD_CMD := \
		DOCKER_BUILDKIT=1 \
		SOURCE_DATE_EPOCH=1 \
		$(BUILDER) \
			build \
			--ulimit nofile=2048:16384 \
			-t $(REGISTRY)/$(NAME):$(VERSION) \
			--build-arg REGISTRY=$(REGISTRY) \
			--build-arg LABEL=$(NAME) \
			--platform $(PLATFORM) \
			$(if $(DOCKER_CACHE_SRC),$(DOCKER_CACHE_SRC),) \
			$(if $(DOCKER_CACHE_DST),$(DOCKER_CACHE_DST),) \
			--network=host \
			--progress=plain \
			$(if $(filter latest,$(VERSION)),,--build-arg VERSION=$(VERSION)) \
			--target $(NAME) \
			-f $(SRC_DIR)/Dockerfile \
			$(EXTRA_ARGS) \
			. \
	)
	$(eval TIMESTAMP := $(shell TZ=GMT date +"%Y-%m-%dT%H:%M:%SZ"))
	mkdir -p out/
	echo $(TIMESTAMP) $(BUILD_CMD) >> out/build.log
	$(BUILD_CMD)
	$(if $(filter package,$(TARGET)),$(BUILDER) save $(REGISTRY)/$(NAME):$(VERSION) -o $@.docker.tar,)
endef

define determine_platform
    ifeq ($1,linux-x86_64)
        PLATFORM := linux/amd64
    else ifeq ($1,linux-aarch64)
        PLATFORM := linux/arm64
    else ifeq ($1,darwin-x86_64)
        PLATFORM := linux/amd64
    else ifeq ($1,darwin-aarch64)
        PLATFORM := linux/arm64
    endif
endef

define go-build
	$(call build,$(3),$(2),latest)
	# Ignore errors from the docker rm; this is just to ensure no such container exists before we create it.
	docker rm -f $(2) 2> /dev/null
	docker create --name=$(2) local/$(2)
	docker export $(2) -o $(OUT_DIR)/$(2).tar
	tar xf $(OUT_DIR)/$(2).tar -C $(OUT_DIR) app
	mv $(OUT_DIR)/app $(OUT_DIR)/$(2)
endef

$(OUT_DIR)/turnkey.%:
	$(call go-build,cmd/turnkey,turnkey,$*)
	mv $(OUT_DIR)/turnkey $@

.DEFAULT_GOAL :=
.PHONY: default
default: \
	$(DEFAULT_GOAL) \
	$(patsubst %,$(KEY_DIR)/%.asc,$(KEYS)) \
	$(OUT_DIR)/turnkey.linux-x86_64 \
	$(OUT_DIR)/turnkey.linux-aarch64 \
	$(OUT_DIR)/turnkey.darwin-x86_64 \
	$(OUT_DIR)/turnkey.darwin-aarch64 \
	$(OUT_DIR)/Formula/turnkey.rb \
	$(OUT_DIR)/release.env \
	$(OUT_DIR)/manifest.txt

.PHONY: lint
lint:
	echo "Running lint"; \
	cd $(SRC_DIR); \
	golangci-lint run ./cmd/turnkey/... --timeout=3m || exit 1;

.PHONY: test
test: $(OUT_DIR)/turnkey.linux-x86_64
	echo "Running tests..."; \
	cd $(SRC_DIR); \
	go test -v ./cmd/turnkey/...

.PHONY: install
install: default
	mkdir -p ~/.local/bin
	cp $(OUT_DIR)/turnkey.$(HOST_OS)-$(HOST_ARCH) ~/.local/bin/turnkey

# Clean repo back to initial clone state
.PHONY: clean
clean:
	git clean -dfx $(SRC_DIR)
	rm -rf $(LOCAL_BUILD_DIR)

$(KEY_DIR)/%.asc:
	$(call fetch_pgp_key,$(basename $(notdir $@)))

$(OUT_DIR)/Formula/turnkey.rb: \
	$(OUT_DIR)/turnkey.darwin-x86_64 \
	$(OUT_DIR)/turnkey.darwin-aarch64
	mkdir -p $(OUT_DIR)/Formula
	export \
		VERSION="$(VERSION)" \
		DARWIN_X86_64_SHA256="$(shell \
			openssl sha256 -r $(OUT_DIR)/turnkey.darwin-x86_64 \
			| sed -e 's/ \*out\// /g' -e 's/ \.\// /g' -e 's/ .*//g' \
		)" \
		DARWIN_AARCH64_SHA256="$(shell \
			openssl sha256 -r $(OUT_DIR)/turnkey.darwin-aarch64 \
			| sed -e 's/ \*out\// /g' -e 's/ \.\// /g' -e 's/ .*//g' \
		)"; \
	cat $(SRC_DIR)/brew/formula.rb | envsubst > $@

$(OUT_DIR)/release.env: | $(OUT_DIR)
	echo 'VERSION=$(VERSION)'              > $(OUT_DIR)/release.env
	echo 'GIT_REF=$(GIT_REF)'             >> $(OUT_DIR)/release.env
	echo 'GIT_AUTHOR=$(GIT_AUTHOR)'       >> $(OUT_DIR)/release.env
	echo 'GIT_KEY=$(GIT_KEY)'             >> $(OUT_DIR)/release.env
	echo 'GIT_TIMESTAMP=$(GIT_TIMESTAMP)' >> $(OUT_DIR)/release.env

.PHONY: build-local
build-local:
	pushd $(shell git rev-parse --show-toplevel)/src/cmd/turnkey; \
	go build -o ../$(LOCAL_BUILD_DIR)/turnkey; \
	popd;