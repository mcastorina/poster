.DEFAULT_TARGET := build

GO       := go
GOFLAGS  := CGO_ENABLED=1
GOFILES  := $(shell find cmd internal -type f -name '*.go' -not -name '*_test.go')

.phony: build
build: bin/poster

.phony: clean
clean: $(shell find bin -type f 2>/dev/null)
	$(if $^, rm -f $^)

bin/poster: cmd/main.go $(GOFILES)
	$(GOFLAGS) $(GO) build -o $@ ./$(<D)
