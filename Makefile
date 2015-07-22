NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m
DEPS=$(shell go list -f '{{range .TestImports}}{{.}} {{end}}' ./...)
GOPKGS=$(shell go list -f '{{.ImportPath}}' ./...)
PKGSDIRS=$(shell go list -f '{{.Dir}}' ./...)
VERSION=$(shell echo `whoami`-`git rev-parse --short HEAD`-`date -u +%Y%m%d%H%M%S`)
DIST_FIND_BUILDS=find . -type d -mindepth 1 -exec

.PHONY: all dist format lint vet build test setup tools deps updatedeps bench clean
.SILENT: all dist format lint vet build test setup tools deps updatedeps bench clean

all: clean build dist

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
	export CGOENABLED=1 && \
	export GOPATH=$(shell godep path):$(shell echo $$GOPATH) &&\
	gox -verbose \
	-ldflags="-X main.version \
	$(VERSION)" \
	-os="windows linux darwin " \
	-arch="amd64" \
	-output="build/{{.OS}}/{{.Dir}}" ./...

dist: build
	@echo "$(OK_COLOR)==> Distro'ing$(NO_COLOR)"
	cd build && \
	$(DIST_FIND_BUILDS) cp ../LICENSE.md {} \; && \
	$(DIST_FIND_BUILDS) cp ../README.md {}/README \; && \
	$(DIST_FIND_BUILDS) cp ../kubesetup.yml {} \; && \
	$(DIST_FIND_BUILDS) zip -qr {}-hpcloud-kubesetup {} \; && \
	cd ..

test: #deps
	@echo "$(OK_COLOR)==> Testing$(NO_COLOR)"
	godep go test -ldflags -linkmode=external -covermode=count ./...

setup:
	@echo "$(OK_COLOR)==> Running initial setup$(NO_COLOR)"
	gox -build-toolchain

tools:
	@echo "$(OK_COLOR)==> Installing tools$(NO_COLOR)"
	#Great tools to have, and used in the build file
	go get -u golang.org/x/tools/cmd/goimports
	go get -u golang.org/x/tools/cmd/vet
	go get -u golang.org/x/tools/cmd/cover
	go get -u golang.org/x/tools/cmd/oracle
	go get -u github.com/golang/lint/golint
	go get -u golang.org/x/tools/...
	#Tools for the project
	go get -u github.com/tools/godep
	go get -u github.com/mitchellh/gox

bench:
	@echo "$(OK_COLOR)==> Benchmark Testing$(NO_COLOR)"
	godep go test -ldflags -linkmode=external -bench=`find . \( ! -regex '.*/\..*' \) -type d`

clean:
	@echo "$(OK_COLOR)==> Cleaning$(NO_COLOR)"
	rm -rf build
	rm -rf $(GOPATH)/pkg/*
	rm -f $(GOPATH)/bin/hpcloud-kubesetup
