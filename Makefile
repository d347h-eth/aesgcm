GO = go
GOBIN = $(CURDIR)/build
GOBUILD = $(GO) build
GOTEST = $(GO) test ./...

## build: builds the binary
build:
	@echo "Building..."
	$(GOBUILD) -C ./cmd -o $(GOBIN)/aesgcm
	@echo "Done."

## clean: cleans the go cache and build dir
clean:
	go clean -cache
	rm -rf ./build

## test: run unit tests
test:
	$(GOTEST) --timeout 100s

## help: print commands help
help:	Makefile
	@sed -n 's/^##//p' $<
