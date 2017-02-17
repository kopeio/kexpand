all: gocode

DOCKER_REGISTRY=kopeio
ifndef VERSION
  VERSION := git-$(shell git rev-parse --short HEAD)
endif

gocode:
	GO15VENDOREXPERIMENT=1 go install -ldflags "-X main.BuildVersion=${VERSION}" github.com/kopeio/kexpand

gocode_docker:
	GO15VENDOREXPERIMENT=1 CGO_ENABLED=0 GOOS=linux go install -ldflags "-s -X main.BuildVersion=${VERSION}" -a -installsuffix cgo github.com/kopeio/kexpand

gofmt:
	gofmt -w -s main.go
	gofmt -w -s cmd

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

