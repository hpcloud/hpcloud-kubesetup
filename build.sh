# To install gox
# go get github.com/mitchellh/gox
#
# To initialize your environment run
# gox -build-toolchain
#
gox -osarch="darwin/amd64 linux/amd64 windows/amd64"
