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
DIST_DIR := dist

lc = $(subst A,a,$(subst B,b,$(subst C,c,$(subst D,d,$(subst E,e,$(subst F,f,$(subst G,g,$(subst H,h,$(subst I,i,$(subst J,j,$(subst K,k,$(subst L,l,$(subst M,m,$(subst N,n,$(subst O,o,$(subst P,p,$(subst Q,q,$(subst R,r,$(subst S,s,$(subst T,t,$(subst U,u,$(subst V,v,$(subst W,w,$(subst X,x,$(subst Y,y,$(subst Z,z,$1))))))))))))))))))))))))))
altarch = $(subst x86_64,amd64,$(subst aarch64,arm64,$1))
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
	$(eval $(call determine_platform,$(1)))
	echo PLATFORM is $(PLATFORM)
	echo HOST_ARCH_ALT is $(HOST_ARCH_ALT)
	echo HOST_OS is $(HOST_OS)
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
			--build-arg ARCH=$(ARCH) \
			--build-arg HOST_ARCH=$(HOST_ARCH) \
			--build-arg GOARCH=$(HOST_ARCH_ALT) \
			--build-arg GOOS=$(HOST_OS) \
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
endef

define determine_platform
    ifeq ($1,linux-x86_64)
        PLATFORM := linux/amd64
		HOST_ARCH_ALT := amd64
		HOST_OS := linux
    else ifeq ($1,linux-aarch64)
        PLATFORM := linux/arm64
		HOST_ARCH_ALT := arm64
		HOST_OS := linux
    else ifeq ($1,darwin-x86_64)
        PLATFORM := linux/amd64
		HOST_ARCH_ALT := amd64
		HOST_OS := darwin
    else ifeq ($1,darwin-aarch64)
        PLATFORM := linux/arm64
		HOST_ARCH_ALT := arm64
		HOST_OS := darwin
    endif
endef

define go-build
	$(call build,$(3),$(2),latest)
	# $(if $(filter package,$(TARGET)),$(BUILDER) save $(REGISTRY)/$(NAME):$(VERSION) -o $@.docker.tar,)
	# Ignore errors from the docker rm; this is just to ensure no such container exists before we create it.
	docker rm -f $(2) 2> /dev/null
	docker create --name=$(2) local/$(2)
	docker export $(2) -o $(OUT_DIR)/$(2).tar
	tar xf $(OUT_DIR)/$(2).tar -C $(OUT_DIR) app
	rm $(OUT_DIR)/$(2).tar
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
	digests.txt

.PHONY: digests.txt
digests.txt:
	echo "Building digests.txt"; \
	cd $(OUT_DIR); \
	sha256sum turnkey.* > ../digests.txt

.PHONY: lint
lint:
	echo "Running lint"; \
	cd $(SRC_DIR); \
	golangci-lint run ./cmd/turnkey/... --timeout=3m || exit 1;

.PHONY: test
test: $(OUT_DIR)/turnkey.$(HOST_OS)-$(HOST_ARCH)
	echo "Running tests..."; \
	cd $(SRC_DIR); \
	go test -v ./cmd/turnkey/...

.PHONY: install
install: default
	mkdir -p ~/.local/bin
	cp $(OUT_DIR)/turnkey.$(HOST_OS)-$(HOST_ARCH) ~/.local/bin/turnkey

.PHONY: clean
clean:
	rm -rf $(LOCAL_BUILD_DIR)
	rm -rf $(OUT_DIR)

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

.PHONY: reproduce
reproduce: clean default
	diff digests.txt digests-dist.txt \
	|| echo "Warning: digests.txt and digests-dist.txt differ"
	
.PHONY: $(DIST_DIR)
$(DIST_DIR): clean default
	rm -rf $@/*
	cp digests.txt digests-dist.txt
	cp -R $(OUT_DIR)/* $@/