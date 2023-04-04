include $(PWD)/src/toolchain/Makefile

KEYS := \
	6B61ECD76088748C70590D55E90A401336C8AAA9 \
	A8864A8303994E3A18ACD1760CAB4418C834B102 \
	66039AA59D823C8BD68DB062D3EC673DF9843E7B \
	DE050A451E6FAF94C677B58B9361DEC647A087BD

ifneq ("$(wildcard $(ROOT)/src/toolchain)","")
	clone := $(shell git submodule update --init --recursive)
endif

.DEFAULT_GOAL :=
.PHONY: default
default: \
	toolchain \
	$(DEFAULT_GOAL) \
	$(patsubst %,$(KEY_DIR)/%.asc,$(KEYS)) \
	$(OUT_DIR)/turnkey.linux-x86_64 \
	$(OUT_DIR)/turnkey.linux-aarch64 \
	$(OUT_DIR)/turnkey.darwin-x86_64 \
	$(OUT_DIR)/turnkey.darwin-aarch64 \
	$(OUT_DIR)/release.env \
	$(OUT_DIR)/manifest.txt

.PHONY: install
install: default
	mkdir -p ~/.local/bin
	cp $(OUT_DIR)/turnkey.$(HOST_OS)-$(HOST_ARCH) ~/.local/bin/turnkey

.PHONY: test
test: $(OUT_DIR)/turnkey.linux-x86_64
	$(call toolchain,' \
		GOCACHE=/home/build/$(CACHE_DIR) \
		GOPATH=/home/build/$(CACHE_DIR) \
		env -C $(SRC_DIR) go test -v ./... \
	')

.PHONY: sign
sign: $(DIST_DIR)/manifest.txt
	set -e; \
	git config --get user.signingkey 2>&1 >/dev/null || { \
		echo "Error: git user.signingkey is not defined"; \
		exit 1; \
	}; \
	fingerprint=$$(\
		git config --get user.signingkey \
		| sed 's/.*\([A-Z0-9]\{16\}\).*/\1/g' \
	); \
	gpg --armor \
		--detach-sig  \
		--output $(DIST_DIR)/manifest.$${fingerprint}.asc \
		$(DIST_DIR)/manifest.txt

.PHONY: verify
verify: $(DIST_DIR)/manifest.txt
	set -e; \
	gpg --import $(KEY_DIR)/*; \
	for file in $(DIST_DIR)/manifest.*.asc; do \
		echo "\nVerifying: $${file}\n"; \
		gpg --verify $${file} $(DIST_DIR)/manifest.txt; \
	done;

# Clean repo back to initial clone state
.PHONY: clean
clean: toolchain-clean
	git clean -dfx $(SRC_DIR)

$(KEY_DIR)/%.asc:
	$(call fetch_pgp_key,$(basename $(notdir $@)))

$(OUT_DIR)/turnkey.%:
	$(call toolchain,' \
		GOHOSTOS="linux" \
		GOHOSTARCH="amd64" \
		GOOS="$(word 1,$(subst -, ,$(word 2,$(subst ., ,$@))))" \
		GOARCH="$(call altarch,$(word 2,$(subst -, ,$(word 2,$(subst ., ,$@)))))" \
		GOCACHE=/home/build/$(CACHE_DIR) \
		GOPATH=/home/build/$(CACHE_DIR) \
		env -C $(SRC_DIR) \
		go build -o /home/build/$@ main.go \
	')
