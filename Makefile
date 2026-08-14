BINARY  ?= tg
BINARY_PARSER  ?= md2tg
TOKEN   ?=
TARGET  ?=

LDFLAGS  = -s -w
ifneq ($(TOKEN),)
LDFLAGS += -X 'main.CompiledToken=$(TOKEN)'
endif
ifneq ($(TARGET),)
LDFLAGS += -X 'main.CompiledTarget=$(TARGET)'
endif

GO_BUILD = CGO_ENABLED=0 go build -trimpath -ldflags "$(LDFLAGS)"

FFI_OUTPUT_DIR ?= ffi/koffi/native
FFI_DARWIN_ARM64_PACKAGE_DIR ?= ffi/koffi-native-darwin-arm64
FFI_DARWIN_X64_PACKAGE_DIR ?= ffi/koffi-native-darwin-x64
FFI_GOOS := $(shell go env GOOS)
ifeq ($(FFI_GOOS),darwin)
FFI_LIBRARY_NAME := libtgmarkdown.dylib
else ifeq ($(FFI_GOOS),windows)
FFI_LIBRARY_NAME := tgmarkdown.dll
else
FFI_LIBRARY_NAME := libtgmarkdown.so
endif

PLATFORMS = \
	linux/amd64 \
	linux/arm64 \
	darwin/amd64 \
	darwin/arm64 \
	windows/amd64

install-tg:
	@if [ -z "$(TOKEN)" ]; then echo "用法: make install TOKEN=<bot_token> [TARGET=<chat_id>] [BINARY=tg]"; exit 1; fi
	$(GO_BUILD) -o $(BINARY) ./cmd/tg/
	mv $(BINARY) $(GOPATH)/bin/$(BINARY)
	@echo "已安装: $(GOPATH)/bin/$(BINARY)"

install-md2tg:
	$(GO_BUILD) -o $(BINARY_PARSER) ./cmd/md2tg/
	mv $(BINARY_PARSER) $(GOPATH)/bin/$(BINARY_PARSER)
	@echo "已安装: $(GOPATH)/bin/$(BINARY_PARSER)"


dist-tg:
	@mkdir -p dist
	@$(foreach PLATFORM,$(PLATFORMS), \
		$(eval OS   := $(word 1,$(subst /, ,$(PLATFORM)))) \
		$(eval ARCH := $(word 2,$(subst /, ,$(PLATFORM)))) \
		$(eval EXT  := $(if $(filter windows,$(OS)),.exe,)) \
		$(eval OUT  := dist/$(BINARY)-$(OS)-$(ARCH)$(EXT)) \
		GOOS=$(OS) GOARCH=$(ARCH) $(GO_BUILD) -o $(OUT) ./cmd/tg/ && echo "  $(OUT)" ; \
	)

dist-md2tg:
	@mkdir -p dist
	@$(foreach PLATFORM,$(PLATFORMS), \
		$(eval OS   := $(word 1,$(subst /, ,$(PLATFORM)))) \
		$(eval ARCH := $(word 2,$(subst /, ,$(PLATFORM)))) \
		$(eval EXT  := $(if $(filter windows,$(OS)),.exe,)) \
		$(eval OUT  := dist/$(BINARY_PARSER)-$(OS)-$(ARCH)$(EXT)) \
		GOOS=$(OS) GOARCH=$(ARCH) $(GO_BUILD) -o $(OUT) ./cmd/md2tg/ && echo "  $(OUT)" ; \
	)

ffi-rich-markdown:
	@mkdir -p $(FFI_OUTPUT_DIR)
	CGO_ENABLED=1 go build -trimpath -ldflags "-s -w" -buildmode=c-shared -o $(FFI_OUTPUT_DIR)/$(FFI_LIBRARY_NAME) ./ffi/rich_markdown
	@echo "已构建: $(FFI_OUTPUT_DIR)/$(FFI_LIBRARY_NAME)"

ffi-rich-markdown-packages: ffi-rich-markdown-darwin-arm64 ffi-rich-markdown-darwin-x64

ffi-rich-markdown-darwin-arm64:
	@mkdir -p $(FFI_DARWIN_ARM64_PACKAGE_DIR)
	CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build -trimpath -ldflags "-s -w" -buildmode=c-shared -o $(FFI_DARWIN_ARM64_PACKAGE_DIR)/libtgmarkdown.dylib ./ffi/rich_markdown
	@echo "已构建: $(FFI_DARWIN_ARM64_PACKAGE_DIR)/libtgmarkdown.dylib"

ffi-rich-markdown-darwin-x64:
	@mkdir -p $(FFI_DARWIN_X64_PACKAGE_DIR)
	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -trimpath -ldflags "-s -w" -buildmode=c-shared -o $(FFI_DARWIN_X64_PACKAGE_DIR)/libtgmarkdown.dylib ./ffi/rich_markdown
	@echo "已构建: $(FFI_DARWIN_X64_PACKAGE_DIR)/libtgmarkdown.dylib"

clean:
	rm -rf dist/
