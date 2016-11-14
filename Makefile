# TODO: Move entirely to bazel?
.PHONY: images

all: gocode

DOCKER_REGISTRY=kopeio
ifndef VERSION
  VERSION := git-$(shell git rev-parse --short HEAD)
endif

gocode:
	GO15VENDOREXPERIMENT=1 go install -ldflags "-X main.BuildVersion=${VERSION}" github.com/kopeio/kexpand/cmd/...

gofmt:
	gofmt -w -s cmd

syncdeps:
	rsync -avz _vendor/ vendor/


# --------------------------------------------------
# Docker images

push: images
	docker push ${DOCKER_REGISTRY}/kexpand

images:
	bazel run //images:kexpand ${DOCKER_REGISTRY}/kexpand


# --------------------------------------------------
# Continuous integration targets

ci: images test govet
	echo "Done"

govet:
	go vet \
	  github.com/kopeio/kexpand/cmd/...

test:
	go test github.com/kopeio/kexpand/cmd/...

