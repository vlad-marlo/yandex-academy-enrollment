.DEFAULT_GOAL := build

.PHONY: build
build:
	go build -o server.o ./src/cmd/server/main.go

.PHONY: gen
gen:
	swag fmt
	swag init --d src/cmd/server/,src/internal/controller/http/,src/pkg/model/ -o ./src/docs/
	go generate ./...

.PHONY: test
test:
	go test ./src/... -v -coverpkg=./internal/... -coverprofile=coverage.out

.PHONY: c
c:
	go tool cover -func coverage.out

.PHONY: tc
tc: test c

.PHONY: lines
lines:
	git ls-files | xargs wc -l
