

fmt:
	gofmt -s -w *.go */*.go

test:
	go test -cover -race -v ./...

release: css
	go build -ldflags="-w -s" -o alexandria main.go

build: css
	go build ./...

run: css
	go run main.go

clean:
	go clean -i ./...

css: ./assets/src/main.sass
	sassc --omit-map-comment --style compressed --sass ./assets/src/main.sass ./assets/public/main.css
