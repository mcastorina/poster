# Reference
# https://www.gnu.org/software/make/manual/make.html#Automatic-Variables
.DEFAULT_TARGET := build

GO       := go
GO_FILES := $(shell find cmd internal -type f -name '*.go')

.phony: build
build: bin/poster

.phony: clean
clean: $(shell find bin -type f 2>/dev/null)
	$(if $^, rm -f $^)

bin/poster: cmd/main.go $(GO_FILES)
	$(GO) build -o $@ ./$(<D)
