NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m
DEPS=$(shell go list -f '{{range .TestImports}}{{.}} {{end}}' ./...)
GOPKGS=$(shell go list -f '{{.ImportPath}}' ./...)
PKGSDIRS=$(shell go list -f '{{.Dir}}' ./...)
VERSION=$(shell echo `whoami`-`git rev-parse --short HEAD`-`date -u +%Y%m%d%H%M%S`)

.PHONY: all dist format lint vet build test setup tools deps updatedeps bench clean
.SILENT: all dist format lint vet build test setup tools deps updatedeps bench clean

all: clean format lint vet build test dist

format:
	@echo "$(OK_COLOR)==> Checking format$(ERROR_COLOR)"
	@echo $(PKGSDIRS) | xargs -I '{p}' -n1 goimports -e -l {p} | sed "s/^/Failed: /"
	@echo "$(NO_COLOR)\c"

lint:
	@echo "$(OK_COLOR)==> Linting$(ERROR_COLOR)"
	@echo $(PKGSDIRS) | xargs -I '{p}' -n1 golint {p}  | sed "s/^/Failed: /"
	@echo "$(NO_COLOR)\c"

vet:
	@echo "$(OK_COLOR)==> Vetting$(ERROR_COLOR)"
	@echo $(GOPKGS) | xargs -I '{p}' -n1 go vet {p}  | sed "s/^/Failed: /"
	@echo "$(NO_COLOR)\c"

build: #deps
	@echo "$(OK_COLOR)==> Building$(NO_COLOR)"
	go build ./...

tools:
	@echo "$(OK_COLOR)==> Installing tools$(NO_COLOR)"
	#Great tools to have, and used in the build file
	go get golang.org/x/tools/cmd/goimports
	go get golang.org/x/tools/cmd/vet
	go get golang.org/x/tools/cmd/cover
	go get github.com/golang/lint/golint

test: #deps
	@echo "$(OK_COLOR)==> Testing$(NO_COLOR)"
	go test $(TEST_FLAGS) -short -ldflags -linkmode=external -covermode=count ./...

integrationtest: #deps
	@echo "$(OK_COLOR)==> Testing$(NO_COLOR)"
	go test $(TEST_FLAGS)  -ldflags -linkmode=external -covermode=count ./...

bench:
	@echo "$(OK_COLOR)==> Benchmark Testing$(NO_COLOR)"
	godep go test $(TEST_FLAGS) -ldflags -linkmode=external -bench=`find . \( ! -regex '.*/\..*' \) -type d`

clean:
	@echo "$(OK_COLOR)==> Cleaning$(NO_COLOR)"
	rm -rf build
	rm -rf $(GOPATH)/pkg/*
	rm -f $(GOPATH)/bin/als
