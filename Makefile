all: gocode

DOCKER_REGISTRY=kopeio
ifndef VERSION
  VERSION := git-$(shell git rev-parse --short HEAD)
endif

gocode:
	GO15VENDOREXPERIMENT=1 go install -ldflags "-X main.BuildVersion=${VERSION}" github.com/kopeio/kexpand

gofmt:
	gofmt -w -s main.go
	gofmt -w -s cmd


builder-image:
	docker build -f images/kexpand-builder/Dockerfile -t kexpand-builder .

build-in-docker: builder-image
	docker run -it -v `pwd`:/src kexpand-builder /onbuild.sh

image: build-in-docker
	docker build -t ${DOCKER_REGISTRY}/kexpand  -f images/kexpand/Dockerfile .

push: image
	docker push ${DOCKER_REGISTRY}/kexpand

syncdeps:
	rsync -avz _vendor/ vendor/
