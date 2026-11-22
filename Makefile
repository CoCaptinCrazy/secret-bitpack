.PHONY: build test vet fmt clean

build:
	go build -o bitpack.exe ./cmd/bitpack

test:
	go test ./...

vet:
	go vet ./...

fmt:
	go fmt ./...

clean:
	if exist bitpack.exe del /q bitpack.exe
