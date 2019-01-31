

fmt:
	gofmt -s -w *.go */*.go

test:
	go test -cover -race -v ./...

release:
	go build -ldflags="-w -s" -o alexandria main.go

build:
	go build ./...

run:
	go run main.go

clean:
	go clean -i ./...
