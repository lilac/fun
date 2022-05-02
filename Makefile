EXE := fun

# install required dependencies
dep:
	go install golang.org/x/tools/cmd/goyacc@latest

pkg/syntax/grammar.go: pkg/syntax/grammar.go.y
	go generate ./pkg

build: $(EXE)

$(EXE): pkg/syntax/grammar.go
	go build ./cmd/fun

test: pkg/syntax/grammar.go
	go test ./...

install: build
	cp $(EXE) $(GOPATH)/bin/
