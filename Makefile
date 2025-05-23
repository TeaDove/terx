GO ?= GO111MODULE=on CGO_ENABLED=0 go
VERSION ?= $(shell cat VERSION)

tag:
	git tag utils/$(VERSION)
	git push origin --tags

test:
	$(GO) test ./...

lint:
	gofumpt -w .
	golines --base-formatter=gofumpt --max-len=120 --no-reformat-tags -w .
	wsl --fix ./...
	golangci-lint run --fix

test-and-lint: test lint