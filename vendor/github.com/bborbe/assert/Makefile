default: test

glide:
	go get github.com/Masterminds/glide

test: glide
	GO15VENDOREXPERIMENT=1 go test -cover `glide novendor`

goimports:
	go get golang.org/x/tools/cmd/goimports

format: goimports
	find . -type f -name '*.go' -not -path './vendor/*' -exec gofmt -w "{}" +
	find . -type f -name '*.go' -not -path './vendor/*' -exec goimports -w "{}" +
