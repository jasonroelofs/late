all: build test vet

build:

test:
	go test ./...

vet:
	go vet ./...

docs: FORCE
	go run scripts/docs.go docs/

FORCE: ;
