EXE := fun

# install required dependencies
dep:
	go install golang.org/x/tools/cmd/goyacc@latest

syntax/grammar.go: syntax/grammar.go.y
	go generate

build: $(EXE)

$(EXE): syntax/grammar.go
	go build -o $(EXE)

test: build
	go test ./...

install: build
	cp $(EXE) $(GOPATH)/bin/
