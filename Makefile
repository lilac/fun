# install required dependencies
dep:
	go install github.com/kivikakk/golex@latest
	go install golang.org/x/tools/cmd/goyacc@latest
	go install github.com/blynn/nex@latest

lexer:
	nex -p fun syntax/fun.nex

syntax/grammar.go: syntax/grammar.go.y
	go generate

build: syntax/grammar.go
	go build

test:
	go test ./...
