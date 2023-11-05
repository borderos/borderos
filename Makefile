CMDS := brd
GOOS := linux
LDFLAGS := -s -w

export GOOS

.PHONY: all
all: $(CMDS)

.PHONY: clean
clean:
	-rm -f $(CMDS)

%: ./cmd/%/*.go
	go build -ldflags="$(LDFLAGS) -X=main.buildTime=$$(date +%Y-%m-%dT%H:%M:%S)" ./cmd/$@

.PHONY: mips-all
mips-all:
	@GOARCH=mips64 $(MAKE)

.PHONY: mips-%
mips-%: ./cmd/%/*.go
	@GOARCH=mips64 $(MAKE) $*
