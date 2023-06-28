# This Makefile is meant to be used by people that do not usually work
# with Go source code. If you know what GOPATH is then you probably
# don't need to bother with make.

.PHONY: geth android ios geth-cross evm all test clean
.PHONY: geth-linux geth-linux-386 geth-linux-amd64 geth-linux-mips64 geth-linux-mips64le
.PHONY: geth-linux-arm geth-linux-arm-5 geth-linux-arm-6 geth-linux-arm-7 geth-linux-arm64
.PHONY: geth-darwin geth-darwin-amd64
.PHONY: geth-windows geth-windows-386 geth-windows-amd64
.PHONY: prepare prepare-system-contracts prepare-ethersjs-project

GOBIN = ./build/bin
GO ?= latest
GORUN = env GO111MODULE=on go run

LSB_exists := $(shell command -v lsb_release 2> /dev/null)

OS :=
ifeq ("$(LSB_exists)","")
	OS = darwin
else
	OS = linux
endif

# We checkout the monorepo as a sibling to the celo-blockchain dir because the
# huge amount of files in the monorepo interferes with tooling such as gopls,
# which becomes very slow.
MONOREPO_PATH=../.celo-blockchain-monorepo-checkouts/$(shell cat monorepo_commit)

# This either evaluates to the contract source files if they exist or NOT_FOUND
# if celo-monorepo has not been checked out yet.
CONTRACT_SOURCE_FILES=$(shell 2>/dev/null find $(MONOREPO_PATH)/packages/protocol \
						   -not -path "*/node_modules*" \
						   -not -path "$(MONOREPO_PATH)/packages/protocol/test*" \
						   -not -path "$(MONOREPO_PATH)/packages/protocol/build*" \
						   -not -path "$(MONOREPO_PATH)/packages/protocol/types*" \
						   || echo "NOT_FOUND")

# example NDK values
export NDK_VERSION ?= android-ndk-r19c
export ANDROID_NDK ?= $(PWD)/ndk_bundle/$(NDK_VERSION)

geth:
	$(GORUN) build/ci.go install ./cmd/geth
	@echo "Done building."
	@echo "Run \"$(GOBIN)/geth\" to launch geth."

prepare: prepare-system-contracts prepare-ethersjs-project

prepare-ethersjs-project: ./e2e_test/ethersjs-api-check/node_modules

./e2e_test/ethersjs-api-check/node_modules: ./e2e_test/ethersjs-api-check/package.json ./e2e_test/ethersjs-api-check/package-lock.json
	@cd ./e2e_test/ethersjs-api-check && npm ci

# This rule checks out celo-monorepo under MONOREPO_PATH at the commit contained in
# monorepo_commit and compiles the system solidity contracts. It then copies the
# compiled contracts from the monorepo to the compiled-system-contracts, so
# that this repo can always access the contracts at a consistent path.
prepare-system-contracts: $(MONOREPO_PATH)/packages/protocol/build
	@rm -rf compiled-system-contracts
	@mkdir -p compiled-system-contracts
	@cp -a $(MONOREPO_PATH)/packages/protocol/build/{contracts,contracts-mento}/ compiled-system-contracts

# If any of the source files in CONTRACT_SOURCE_FILES are more recent than the
# build dir or the build dir does not exist then we remove the build dir, yarn
# install and rebuild the contracts.
$(MONOREPO_PATH)/packages/protocol/build: $(CONTRACT_SOURCE_FILES)
	@node --version | grep "^v18" || (echo "node v18 is required to build the monorepo (nvm use 18)" && exit 1)
	@echo Running yarn install and compiling contracts
	@cd $(MONOREPO_PATH) && rm -rf packages/protocol/build && yarn && cd packages/protocol && yarn run build:sol


# This target serves as an intermediate step to avoid running the
# $(MONOREPO_PATH) target once per contract source file.  This could also be
# achieved by using the group targets separator '&:' instead of just ':', but
# that functionality was added in make version 4.3 which doesn't seem to be
# readily available on most systems yet. So although this rule will be run once
# for each source file, since it is empty that is very quick. $(MONOREPO_PATH)
# as a prerequisite of this will be run at most once.
$(CONTRACT_SOURCE_FILES): $(MONOREPO_PATH)

# Clone the monorepo at the commit in the file `monorepo_commit`.
#
# The checkouts are kept separate by commit to make switching between commits. Use `make clean-old-monorepos` to remove all checkouts but the once currently indicated by `monorepo_commit`.
$(MONOREPO_PATH): monorepo_commit
	@set -e; \
	mc=`cat monorepo_commit`; \
	echo "monorepo_commit is $${mc}"; \
	if git ls-remote --heads --exit-code git@github.com:celo-org/celo-monorepo.git $${mc} > /dev/null 2>&1; \
	then \
		echo "Expected commit hash or tag in 'monorepo_commit' instead found branch name '$${mc}'"; \
		exit 1; \
	fi; \
	if  [ ! -e $(MONOREPO_PATH) ]; \
	then \
		echo "Cloning monorepo at $${mc}"; \
		mkdir -p $(MONOREPO_PATH) && cd $(MONOREPO_PATH); \
		git init; \
		git remote add origin https://github.com/celo-org/celo-monorepo.git; \
		git fetch --quiet --depth 1 origin $${mc}; \
		git checkout FETCH_HEAD; \
	fi


clean-old-monorepos:
	@all_repos_dir=$$(realpath $(MONOREPO_PATH)/..); \
	delete_dirs=$$(ls -d -1 $${all_repos_dir}/*/ | grep -v $$(cat monorepo_commit)); \
	echo Deleting $$delete_dirs; \
	rm -fr $$delete_dirs


geth-musl:
	$(GORUN) build/ci.go install -musl ./cmd/geth
	@echo "Done building with musl."
	@echo "Run \"$(GOBIN)/geth\" to launch geth."

check_android_env:
	@test $${ANDROID_NDK?Please set environment variable ANDROID_NDK}
	@test $${ANDROID_HOME?Please set environment variable ANDROID_HOME}

ndk_bundle: check_android_env
ifeq ("$(wildcard $(ANDROID_NDK))","")
	@test $${NDK_VERSION?Please set environment variable NDK_VERSION}
	curl --silent --show-error --location --fail --retry 3 --output /tmp/$(NDK_VERSION).zip \
		https://dl.google.com/android/repository/$(NDK_VERSION)-$(OS)-x86_64.zip && \
		rm -rf $(ANDROID_NDK) && \
		mkdir -p $(ANDROID_NDK) && \
		unzip -q /tmp/$(NDK_VERSION).zip -d $(ANDROID_NDK)/.. && \
		rm /tmp/$(NDK_VERSION).zip
else
ifeq ("$(wildcard $(ANDROID_NDK)/toolchains/llvm/prebuilt/$(OS)-x86_64)","")
	$(error "Android NDK is installed but doesn't contain an llvm cross-compilation toolchain. Delete your current NDK or modify the ANDROID_NDK environment variable to an empty directory download it automatically.")
endif
endif

swarm:
	$(GORUN) build/ci.go install ./cmd/swarm
	@echo "Done building."
	@echo "Run \"$(GOBIN)/swarm\" to launch swarm."

all:
	$(GORUN) build/ci.go install

all-musl:
	$(GORUN) build/ci.go install -musl

android:
	@echo "Applying patch for mobile libs..."
	git apply patches/mobileLibsForBuild.patch
	ANDROID_NDK_HOME=$(ANDROID_NDK) $(GORUN) build/ci.go aar --local --metrics-default
	@echo "Done building."
	@echo "Import \"$(GOBIN)/geth.aar\" to use the library."
	@echo "Import \"$(GOBIN)/geth-sources.jar\" to add javadocs"
	@echo "For more info see https://stackoverflow.com/questions/20994336/android-studio-how-to-attach-javadoc"
	@echo "Remove patch for mobile libs..."
	git apply -R patches/mobileLibsForBuild.patch
	

ios:
	CGO_ENABLED=1 DISABLE_BITCODE=true $(GORUN) build/ci.go xcode --local --metrics-default
	pushd "$(GOBIN)"; rm -rf Geth.framework.tgz; tar -czvf Geth.framework.tgz Geth.framework; popd
	# Geth.framework is a static framework, so we have to also keep the other static libs it depends on
	# in order to link it to the final app
	# One day gomobile will probably support xcframework which would solve this ;-)
	cp -f "$$(go list -m -f "{{ .Dir }}" github.com/celo-org/celo-bls-go-ios)/libs/universal/libbls_snark_sys.a" .
	@echo "Done building."
	@echo "Import \"$(GOBIN)/Geth.framework\" to use the library."

test: all
	$(GORUN) build/ci.go test $(TEST_FLAGS)

lint: ## Run linters.
	$(GORUN) build/ci.go lint

clean-geth:
	env GO111MODULE=on go clean -cache
	rm -fr build/_workspace/pkg/ $(GOBIN)/*

clean: clean-geth

# The devtools target installs tools required for 'go generate'.
# You need to put $GOBIN (or $GOPATH/bin) in your PATH to use 'go generate'.

devtools:
	env GOBIN= go install golang.org/x/tools/cmd/stringer@latest
	env GOBIN= go install github.com/kevinburke/go-bindata/go-bindata@latest
	env GOBIN= go install github.com/fjl/gencodec@latest
	env GOBIN= go install github.com/golang/protobuf/protoc-gen-go@latest
	env GOBIN= go install ./cmd/abigen
	@type "solc" 2> /dev/null || echo 'Please install solc'
	@type "protoc" 2> /dev/null || echo 'Please install protoc'

# Cross Compilation Targets (xgo)

geth-cross: geth-linux geth-darwin geth-windows geth-android geth-ios
	@echo "Full cross compilation done:"
	@ls -ld $(GOBIN)/geth-*

geth-linux: geth-linux-386 geth-linux-amd64 geth-linux-arm geth-linux-mips64 geth-linux-mips64le
	@echo "Linux cross compilation done:"
	@ls -ld $(GOBIN)/geth-linux-*

geth-linux-386:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/386 -v ./cmd/geth
	@echo "Linux 386 cross compilation done:"
	@ls -ld $(GOBIN)/geth-linux-* | grep 386

geth-linux-amd64:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/amd64 -v ./cmd/geth
	@echo "Linux amd64 cross compilation done:"
	@ls -ld $(GOBIN)/geth-linux-* | grep amd64

geth-linux-arm: geth-linux-arm-5 geth-linux-arm-6 geth-linux-arm-7 geth-linux-arm64
	@echo "Linux ARM cross compilation done:"
	@ls -ld $(GOBIN)/geth-linux-* | grep arm

geth-linux-arm-5:
	# requires an arm compiler, on Ubuntu: sudo apt-get install gcc-arm-linux-gnueabi g++-arm-linux-gnueabi
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/arm-5 -v ./cmd/geth
	@echo "Linux ARMv5 cross compilation done:"
	@ls -ld $(GOBIN)/geth-linux-* | grep arm-5

geth-linux-arm-6:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/arm-6 -v ./cmd/geth
	@echo "Linux ARMv6 cross compilation done:"
	@ls -ld $(GOBIN)/geth-linux-* | grep arm-6

geth-linux-arm-7:
	# requires an arm compiler, on Ubuntu: sudo apt-get install gcc-arm-linux-gnueabihf g++-arm-linux-gnueabihf
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/arm-7 -v  --tags arm7 ./cmd/geth
	@echo "Linux ARMv7 cross compilation done:"
	@ls -ld $(GOBIN)/geth-linux-* | grep arm-7

geth-linux-arm64:
	# requires an arm64 compiler, on Ubuntu: sudo apt-get install gcc-aarch64-linux-gnu	g++-aarch64-linux-gnu
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/arm64 -v ./cmd/geth
	@echo "Linux ARM64 cross compilation done:"
	@ls -ld $(GOBIN)/geth-linux-* | grep arm64

geth-linux-mips:
	# requires a mips compiler, on Ubuntu: sudo apt-get install gcc-mips-linux-gnu
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/mips --ldflags '-extldflags "-static"' -v ./cmd/geth
	@echo "Linux MIPS cross compilation done:"
	@ls -ld $(GOBIN)/geth-linux-* | grep mips

geth-linux-mipsle:
	# requires a mips compiler, on Ubuntu: sudo apt-get install gcc-mipsel-linux-gnu
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/mipsle --ldflags '-extldflags "-static"' -v ./cmd/geth
	@echo "Linux MIPSle cross compilation done:"
	@ls -ld $(GOBIN)/geth-linux-* | grep mipsle

geth-linux-mips64:
	# requires a mips compiler, on Ubuntu: sudo apt-get install gcc-mips64-linux-gnuabi64
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/mips64 --ldflags '-extldflags "-static"' -v ./cmd/geth
	@echo "Linux MIPS64 cross compilation done:"
	@ls -ld $(GOBIN)/geth-linux-* | grep mips64

geth-linux-mips64le:
	# requires a mips compiler, on Ubuntu: sudo apt-get install gcc-mips64el-linux-gnuabi64
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/mips64le --ldflags '-extldflags "-static"' -v ./cmd/geth
	@echo "Linux MIPS64le cross compilation done:"
	@ls -ld $(GOBIN)/geth-linux-* | grep mips64le

geth-darwin: geth-darwin-amd64
	@echo "Darwin cross compilation done:"
	@ls -ld $(GOBIN)/geth-darwin-*

geth-darwin-amd64:
	# needs include files for asm errno, on Ubuntu: sudo apt-get install linux-libc-dev
	# currently doesn't compile on Ubuntu
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=darwin/amd64 -v ./cmd/geth
	@echo "Darwin amd64 cross compilation done:"
	@ls -ld $(GOBIN)/geth-darwin-* | grep amd64

geth-windows: geth-windows-386 geth-windows-amd64
	@echo "Windows cross compilation done:"
	@ls -ld $(GOBIN)/geth-windows-*

geth-windows-386:
	# currently doesn't compile on Ubuntu, missing libunwind in xgo
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=windows/386 -v ./cmd/geth
	@echo "Windows 386 cross compilation done:"
	@ls -ld $(GOBIN)/geth-windows-* | grep 386

geth-windows-amd64:
	# currently doesn't compile on Ubuntu, missing libunwind in xgo
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=windows/amd64 -v ./cmd/geth
	@echo "Windows amd64 cross compilation done:"
	@ls -ld $(GOBIN)/geth-windows-* | grep amd64
