.DEFAULT_GOAL := build

.PHONY: build
build:
	go build -o server ./cmd/server/main.go

.PHONY: gen
gen:
	go generate ./...

.PHONY: test
test:
	go test ./... -v -coverpkg=./internal/... -coverprofile=coverage.out

.PHONY: c
c:
	go tool cover -func coverate.out

.PHONY: lines
lines:
	git ls-files | xargs wc -l
