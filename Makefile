# install required dependencies
dep:
	go install golang.org/x/tools/cmd/goyacc@latest

syntax/grammar.go: syntax/grammar.go.y
	go generate

build: syntax/grammar.go
	go build

test: build
	go test ./...
