all: test vet

test:
	go test ./...

vet:
	go vet ./...
