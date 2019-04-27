.PHONY: build
build: css
	go build ./...

.PHONY: fmt
fmt:
	gofmt -s -w *.go */*.go

.PHONY: vet
vet:
	go vet ./...

.PHONY: test
test: vet
	go test -cover -race -v ./...

.PHONY: lint
lint:
	golint ./...

.PHONY: release
release: css
	go build -ldflags="-w -s" -o alexandria main.go

.PHONY: run
run: css
	go run main.go

.PHONY: clean
clean:
	go clean -i ./...

.PHONY: css
css: ./assets/src/main.sass
	sassc --omit-map-comment --style compressed --sass ./assets/src/main.sass ./assets/public/main.css
