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

.PHONY: build install clean

build:
	go build -ldflags "$(LDFLAGS)" -o $(BINARY) .

# 锁定 token/target 并安装到 PATH
install:
	@if [ -z "$(TOKEN)" ]; then echo "用法: make install TOKEN=<bot_token> [TARGET=<chat_id>] [BINARY=tg]"; exit 1; fi
	go build -ldflags "$(LDFLAGS)" -o $(BINARY) .
	mv $(BINARY) $(GOPATH)/bin/$(BINARY)
	@echo "已安装: $(GOPATH)/bin/$(BINARY)"

clean:
	rm -f $(BINARY)
