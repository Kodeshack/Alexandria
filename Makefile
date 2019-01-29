

fmt:
	gofmt -s -w *.go */*.go

test:
	go test -cover -race -v ./...

build:
	go build ./...

run:
	go run main.go

clean:
	go clean -i ./...
