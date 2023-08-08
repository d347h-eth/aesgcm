GO = go
GOBIN = $(CURDIR)/build
GOBUILD = $(GO) build

## build: builds the binary
build:
	@echo "Building..."
	$(GOBUILD) -C ./cmd -o $(GOBIN)/aesgcm
	@echo "Done."

## clean: cleans the go cache and build dir
clean:
	go clean -cache
	rm -rf ./build

## help: print commands help
help:	Makefile
	@sed -n 's/^##//p' $<
