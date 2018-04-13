IMAGE ?= bborbe/k8s-deploy
REGISTRY ?= docker.io
ifeq ($(VERSION),)
	VERSION := $(shell git describe --tags `git rev-list --tags --max-count=1`)
endif

all: test install run

prepare:
	go get -u golang.org/x/tools/cmd/goimports
	go get -u github.com/golang/lint/golint
	go get -u github.com/kisielk/errcheck
	go get -u github.com/golang/dep/cmd/dep
	go get -u github.com/Masterminds/glide
	go get -u github.com/bborbe/docker_utils/bin/docker_remote_tag_exists

glide:
	go get github.com/Masterminds/glide

test: glide
	GO15VENDOREXPERIMENT=1 go test -cover `glide novendor`

install:
	GOBIN=$(GOPATH)/bin GO15VENDOREXPERIMENT=1 go install cmd/k8s-deploy/*.go

run:
	k8s-deploy \
	-dir=~/Documents/workspaces/kubernetes/manifests/netcup/template \
	-namespace=debug \
	-teamvault-config=~/.teamvault.json \
	-logtostderr \
	-v=4

goimports:
	go get golang.org/x/tools/cmd/goimports

format: goimports
	find . -type f -name '*.go' -not -path './vendor/*' -exec gofmt -w "{}" +
	find . -type f -name '*.go' -not -path './vendor/*' -exec goimports -w "{}" +

buildgo:
	CGO_ENABLED=0 GOOS=linux go build -ldflags "-s" -a -installsuffix cgo -o k8s-deploy ./go/src/github.com/seibert-media/k8s-deploy/cmd/k8s-deploy

build:
	docker build --no-cache --rm=true -t $(REGISTRY)/$(IMAGE)-build:$(VERSION) -f ./Dockerfile.build .
	docker run -t $(REGISTRY)/$(IMAGE)-build:$(VERSION) /bin/true
	docker cp `docker ps -q -n=1 -f ancestor=$(REGISTRY)/$(IMAGE)-build:$(VERSION) -f status=exited`:/k8s-deploy .
	docker rm `docker ps -q -n=1 -f ancestor=$(REGISTRY)/$(IMAGE)-build:$(VERSION) -f status=exited` || true
	docker build --no-cache --rm=true --tag=$(REGISTRY)/$(IMAGE):$(VERSION) -f Dockerfile.static .
	rm -f k8s-deploy

upload:
	docker push $(REGISTRY)/$(IMAGE):$(VERSION)

clean:
	docker rmi $(REGISTRY)/$(IMAGE):$(VERSION) || true

version:
	@echo $(VERSION)

docker_remote_tag_exists:
	@go get github.com/bborbe/docker_utils/bin/docker_remote_tag_exists

trigger: docker_remote_tag_exists
	@exists=`docker_remote_tag_exists \
		-registry=${REGISTRY} \
		-repository="${IMAGE}" \
		-credentialsfromfile \
		-tag="${VERSION}" \
		-alsologtostderr \
		-v=0`; \
	trigger="build"; \
	if [ "$${exists}" = "true" ]; then \
		trigger="skip"; \
	fi; \
	echo $${trigger}
