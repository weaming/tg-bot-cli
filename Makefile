BINARY  ?= tg
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

.PHONY: build install dist clean

build:
	$(GO_BUILD) -o $(BINARY) .

install:
	@if [ -z "$(TOKEN)" ]; then echo "用法: make install TOKEN=<bot_token> [TARGET=<chat_id>] [BINARY=tg]"; exit 1; fi
	$(GO_BUILD) -o $(BINARY) .
	mv $(BINARY) $(GOPATH)/bin/$(BINARY)
	@echo "已安装: $(GOPATH)/bin/$(BINARY)"

dist:
	@mkdir -p dist
	@$(foreach PLATFORM,$(PLATFORMS), \
		$(eval OS   := $(word 1,$(subst /, ,$(PLATFORM)))) \
		$(eval ARCH := $(word 2,$(subst /, ,$(PLATFORM)))) \
		$(eval EXT  := $(if $(filter windows,$(OS)),.exe,)) \
		$(eval OUT  := dist/$(BINARY)-$(OS)-$(ARCH)$(EXT)) \
		GOOS=$(OS) GOARCH=$(ARCH) $(GO_BUILD) -o $(OUT) . && echo "  $(OUT)" ; \
	)

clean:
	rm -f $(BINARY)
	rm -rf dist/
