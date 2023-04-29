BINARY_NAME=beango

build:
	go build -o ${BINARY_NAME} main.go

run: build
	./${BINARY_NAME}

test-unit:
	go test ./...

testcov:
	go test ./... -coverprofile=coverage.out

dep:
	go mod download

vet:
	go vet

clean:
	go clean
