APPNAME=$(shell basename $(CURDIR))

default:
	go build -o bin/$(APPNAME)

deps:
	go get github.com/mitchellh/gox

publish: check-env deps
	@sh -c "'scripts/release.sh'"

check-env:
	@if test "$(BIN_BUCKET_NAME)" = "" ; then \
		echo "BIN_BUCKET_NAME not set"; \
		exit 1; \
	fi

fmt:
	gofmt -w .

.PHONY: default
