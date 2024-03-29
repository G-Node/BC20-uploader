# Binary
APP = uploader

# Build loc
BUILDLOC = build

# Install location
INSTLOC = $(GOPATH)/bin

cwd = $(shell pwd)

# Build flags
ncommits = $(shell git rev-list --count HEAD)
BUILDNUM = $(shell printf '%06d' $(ncommits))
COMMITHASH = $(shell git rev-parse HEAD)
LDFLAGS = -ldflags="-X main.build=$(BUILDNUM) -X main.commit=$(COMMITHASH)"

SOURCES = $(shell find . -type f -iname "*.go") go.mod go.sum

.PHONY: $(APP) clean uninstall

$(APP): $(BUILDLOC)/$(APP)

clean:
	rm -rf $(BUILDLOC)

uninstall:
	rm $(INSTLOC)/$(APP)

$(BUILDLOC)/$(APP): $(SOURCES)
	go build -trimpath $(LDFLAGS) $(GCFLAGS) -o $(BUILDLOC)/$(APP) ./cmd/uploader
