include $(PWD)/src/toolchain/Makefile

ifneq ("$(wildcard $(ROOT)/src/toolchain)","")
	clone := $(shell git submodule update --init --recursive)
endif

.DEFAULT_GOAL :=
.PHONY: default
default: \
	toolchain \
	$(DEFAULT_GOAL) \
	$(OUT_DIR)/turnkey.linux-386 \
	$(OUT_DIR)/turnkey.linux-amd64 \
	$(OUT_DIR)/turnkey.linux-arm64 \
	$(OUT_DIR)/turnkey.darwin-amd64 \
	$(OUT_DIR)/turnkey.darwin-arm64 \
	$(OUT_DIR)/release.env \
	$(OUT_DIR)/manifest.txt

.PHONY: test
test: $(OUT_DIR)/turnkey.linux-amd64
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
	for file in $(DIST_DIR)/manifest.*.asc; do \
		echo "\nVerifying: $${file}\n"; \
		gpg --verify $${file} $(DIST_DIR)/manifest.txt; \
	done;

# Clean repo back to initial clone state
.PHONY: clean
clean: toolchain-clean
	git clean -dfx $(SRC_DIR)

$(OUT_DIR)/turnkey.%:
	$(call toolchain,' \
		GOHOSTOS="linux" \
		GOHOSTARCH="amd64" \
		GOOS="$(word 1,$(subst -, ,$(word 2,$(subst ., ,$@))))" \
		GOARCH="$(word 2,$(subst -, ,$(word 2,$(subst ., ,$@))))" \
		GOCACHE=/home/build/$(CACHE_DIR) \
		GOPATH=/home/build/$(CACHE_DIR) \
		env -C $(SRC_DIR) \
		go build -o /home/build/$@ main.go \
	')
