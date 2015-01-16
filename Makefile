.PHONY: build run test clean

BINARY := entity-extractor
BUILDFILES := config.go extractor.go extractor_api.go main.go
IMPORT_BASE := github.com/alphagov
IMPORT_PATH := $(IMPORT_BASE)/entity-extractor

build: _vendor
	gom build -o $(BINARY) $(BUILDFILES)

run: _vendor
	gom run $(BUILDFILES)

test: _vendor build
	gom test -v

clean:
	rm -f $(BINARY)

_vendor: Gomfile _vendor/src/$(IMPORT_PATH)
	gom install
	touch _vendor

_vendor/src/$(IMPORT_PATH):
	rm -f _vendor/src/$(IMPORT_PATH)
	mkdir -p _vendor/src/$(IMPORT_BASE)
	ln -s $(CURDIR) _vendor/src/$(IMPORT_PATH)
