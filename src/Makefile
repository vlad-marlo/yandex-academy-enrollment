.DEFAULT_GOAL := build

.PHONY: build
build:
	go build -o server.o ./cmd/server/main.go

.PHONY: gen
gen:
	swag fmt
	swag init --d cmd/server/,internal/controller/http/,pkg/model/
	go generate ./...

.PHONY: test
test:
	go test ./... -v -coverpkg=./internal/...,./pkg/... -coverprofile=coverage.out

.PHONY: c
c:
	go tool cover -func coverage.out

.PHONY: tc
tc: test c

.PHONY: lines
lines:
	git ls-files | xargs wc -l
