# Reference
# https://www.gnu.org/software/make/manual/make.html#Automatic-Variables

GO_FILES := $(shell find cmd internal -type f -name '*.go')


all: build

build: bin/poster

clean: $(shell find bin -type f 2>/dev/null)
	$(if $^, rm -f $^)

bin/poster: cmd/main.go $(GO_FILES)
	go build -o $@ ./$(<D)

.phony: run build clean
