NAME=bob
BINARY=terraform-provider-bob
VERSION=0.1.0
OS_ARCH=linux_amd64

default: install

build:
	go build -o ${BINARY}

install: build
	mkdir -p ~/.terraform.d/plugins/bob-cd/providers/bob/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/bob-cd/providers/bob/${VERSION}/${OS_ARCH}

test:
	go test || exit 1
