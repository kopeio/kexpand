all: gocode

DOCKER_REGISTRY=kopeio
UNIQUE:=$(shell date +%s)
GOVERSION=1.7.4
GITSHA := $(git describe --always)
ifndef VERSION
  VERSION=0.2
endif

# See http://stackoverflow.com/questions/18136918/how-to-get-current-relative-directory-of-your-makefile
MAKEDIR:=$(strip $(shell dirname "$(realpath $(lastword $(MAKEFILE_LIST)))"))

gocode:
	GO15VENDOREXPERIMENT=1 go install -ldflags "-X main.BuildVersion=${VERSION}" github.com/kopeio/kexpand

gocode_docker:
	GO15VENDOREXPERIMENT=1 CGO_ENABLED=0 GOOS=linux go install -ldflags "-s -X main.BuildVersion=${VERSION}" -a -installsuffix cgo github.com/kopeio/kexpand

crossbuild-in-docker:
	docker pull golang:${GOVERSION} # Keep golang image up to date
	docker run --name=kexpand-build-${UNIQUE} -e STATIC_BUILD=yes -e VERSION=${VERSION} -v ${MAKEDIR}:/go/src/github.com/kopeio/kexpand golang:${GOVERSION} make -f /go/src/github.com/kopeio/kexpand/Makefile crossbuild
	docker cp kexpand-build-${UNIQUE}:/go/.build .

crossbuild:
	mkdir -p .build/dist/
	GOOS=darwin GOARCH=amd64 go build -a ${EXTRA_BUILDFLAGS} -o .build/dist/darwin/amd64/kexpand -ldflags "${EXTRA_LDFLAGS} -X main.BuildVersion=${VERSION} -X main.GitVersion=${GITSHA}" github.com/kopeio/kexpand
	GOOS=linux GOARCH=amd64 go build -a ${EXTRA_BUILDFLAGS} -o .build/dist/linux/amd64/kexpand -ldflags "${EXTRA_LDFLAGS} -X main.BuildVersion=${VERSION} -X main.GitVersion=${GITSHA}" github.com/kopeio/kexpand


kexpand-dist: crossbuild-in-docker
	mkdir -p .build/dist/
	(sha1sum .build/dist/darwin/amd64/kexpand | cut -d' ' -f1) > .build/dist/darwin/amd64/kexpand.sha1
	(sha1sum .build/dist/linux/amd64/kexpand | cut -d' ' -f1) > .build/dist/linux/amd64/kexpand.sha1

version-dist: kexpand-dist

gofmt:
	gofmt -w -s main.go
	gofmt -w -s cmd
	gofmt -w -s pkg

build-in-docker:
	docker run -it -v `pwd`:/src golang:1.7 /src/images/kexpand/onbuild.sh

image: build-in-docker
	docker build -t ${DOCKER_REGISTRY}/kexpand  -f images/kexpand/Dockerfile .

push: image
	docker push ${DOCKER_REGISTRY}/kexpand


# --------------------------------------------------
# Continuous integration targets

ci: images test govet
	echo "Done"

govet:
	go vet \
	  github.com/kopeio/kexpand/cmd/...

test:
	go test github.com/kopeio/kexpand/cmd/...

