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

dist-tg:
	@mkdir -p dist
	@$(foreach PLATFORM,$(PLATFORMS), \
		$(eval OS   := $(word 1,$(subst /, ,$(PLATFORM)))) \
		$(eval ARCH := $(word 2,$(subst /, ,$(PLATFORM)))) \
		$(eval EXT  := $(if $(filter windows,$(OS)),.exe,)) \
		$(eval OUT  := dist/$(BINARY)-$(OS)-$(ARCH)$(EXT)) \
		GOOS=$(OS) GOARCH=$(ARCH) $(GO_BUILD) -o $(OUT) ./cmd/tg/ && echo "  $(OUT)" ; \
	)

install-md2tg:
	$(GO_BUILD) -o $(BINARY_PARSER) ./cmd/md2tg/
	mv $(BINARY_PARSER) $(GOPATH)/bin/$(BINARY_PARSER)
	@echo "已安装: $(GOPATH)/bin/$(BINARY_PARSER)"

clean:
	rm -f $(BINARY)
	rm -rf dist/
