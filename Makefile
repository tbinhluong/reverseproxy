MAIN_FILE=main.go
BINARY_NAME=reverseproxy
BINARY_LINUX=$(BINARY_NAME)_linux
BINARY_WINDOWS=$(BINARY_NAME)_windows
BINARY_BSD=$(BINARY_NAME)_freebsd
BINARY_DARWIN=$(BINARY_NAME)_darwin

GO=go
GOBUILD=$(GO) build
GOCLEAN=$(GO) clean

DIST_DIR=dist

build:
	$(GOBUILD) -ldflags "-s -w" -o $(DIST_DIR)/$(BINARY_NAME) $(MAIN_FILE)

build_linux:
	GOOS=linux GOARCH=amd64 $(GOBUILD) -ldflags "-s -w" -o $(DIST_DIR)/$(BINARY_LINUX)  $(MAIN_FILE)

build_windows:
	GOOS=windows GOARCH=amd64 $(GOBUILD) -ldflags "-s -w" -o $(DIST_DIR)/$(BINARY_WINDOWS)  $(MAIN_FILE)

build_freebsd:
	GOOS=freebsd GOARCH=amd64 $(GOBUILD) -ldflags "-s -w" -o $(DIST_DIR)/$(BINARY_BSD)  $(MAIN_FILE)

build_darwin:
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -ldflags "-s -w" -o $(DIST_DIR)/$(BINARY_DARWIN)  $(MAIN_FILE)

clean:
	$(GOCLEAN)
	rm -f $(DIST_DIR)/$(BINARY_NAME)
	rm -f $(DIST_DIR)/$(BINARY_LINUX)
	rm -f $(DIST_DIR)/$(BINARY_WINDOWS)
	rm -f $(DIST_DIR)/$(BINARY_BSD)