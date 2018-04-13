default: test

prepare:
	go get -u golang.org/x/tools/cmd/goimports
	go get -u github.com/Masterminds/glide
	go get -u github.com/golang/lint/golint
	go get -u github.com/kisielk/errcheck

check: lint vet errcheck

format: goimports
	find . -type f -name '*.go' -not -path './vendor/*' -exec gofmt -w "{}" +
	go get golang.org/x/tools/cmd/goimports
	find . -type f -name '*.go' -not -path './vendor/*' -exec goimports -w "{}" +

test: glide
	GO15VENDOREXPERIMENT=1 go test -cover `glide novendor`

glide:
	go get github.com/Masterminds/glide

dep: glide
	glide up

vet:
	go tool vet .
	go tool vet --shadow .

lint:
	go get github.com/golang/lint/golint
	golint -min_confidence 1 ./...

errcheck:
	go get github.com/kisielk/errcheck
	errcheck -ignore '(Close|Write)' ./...

