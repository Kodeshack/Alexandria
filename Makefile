

fmt:
	gofmt -s -w *.go */*.go

.PHONY: vet
vet:
	go vet ./...

test: vet
	go test -cover -race -v ./...

release: css codemirror
	go build -ldflags="-w -s" -o alexandria main.go

build: css
	go build ./...

run: css
	go run main.go

clean:
	go clean -i ./...

css: ./assets/src/main.sass
	sassc --omit-map-comment --style compressed --sass ./assets/src/main.sass ./assets/public/main.css

codemirror: ./assets/src/CodeMirror/{lib/codemirror.*,mode/*/*.js,theme/*.css,addon/*.js,keymap/*.js}  ./assets/src/editor.js
	cp ./assets/src/CodeMirror/lib/codemirror.* ./assets/public/
	cp ./assets/src/CodeMirror/mode/markdown/markdown.js ./assets/public/
	cp ./assets/src/CodeMirror/theme/xq-light.css ./assets/public/
	cp ./assets/src/CodeMirror/keymap/vim.js ./assets/public/
	cat ./assets/src/editor.js | tr '\n' ' ' | sed -E -e "s/ {2,}//g" > ./assets/public/editor.js
