
# NOTE: This Makefile is only necessary if you 
# plan on developing the msgp tool and library.
# Installation can still be performed with a
# normal `go install`.


# generated integration test files
GGEN = ./_generated/generated.go ./_generated/generated_test.go
# generated unit test files

MGEN = ./msgp/defgen_test.go ./msgp/nestedgen_test.go

# generated green layer above msgp
ZGEN = ./green/green_gen.go

SHELL := /bin/bash

BIN = $(GOBIN)/greenpack2

.PHONY: clean wipe install get-deps bench all dev

dev: clean install test

$(BIN): */*.go *.go
	@go install && (cd ./cmd/addzid && make)

install:
	/bin/echo "package main" > gitcommit.go
	/bin/echo "func init() { LAST_GIT_COMMIT_HASH = \"$(shell git rev-parse HEAD)\"; NEAREST_GIT_TAG= \"$(shell git describe --abbrev=0 --tags)\"; GIT_BRANCH=\"$(shell git rev-parse --abbrev-ref  HEAD)\"; GO_VERSION=\"$(shell go version)\";}" >> gitcommit.go
	go install  && (cd ./cmd/addzid && make)

$(GGEN): ./_generated/def.go
	go generate ./_generated

$(MGEN): ./msgp/defs_test.go
	go generate ./msgp

$(ZGEN): ./green/green.go
	go install
	go generate ./green

test: all
	go test -v ./parse
	go test -v ./msgp
	go test -v ./_generated
	go test -v ./green
	# and test addzid
	go test -v ./cmd/addzid
	# build and run on testdata/
	go build -o ./greenpack2
	cd testdata && go generate && go test -v
	./greenpack2 -file testdata/my.go && go test -v ./testdata/my_gen_test.go ./testdata/my.go ./testdata/my_gen.go
	./greenpack2 -file testdata/my.go -o testdata/my_msgp_gen.go -method-prefix=MSGP -tests=false -io=false # test the -method-prefix flag


bench: all
	go test -bench . ./msgp
	go test -bench . ./_generated

clean:
	$(RM) $(GGEN) $(MGEN) ./greenpack2

wipe: clean
	$(RM) $(BIN)

get-deps:
	go get -d -t ./...

all: install $(GGEN) $(MGEN) $(ZGEN)

# travis CI enters here
travis:
	go get -d -t ./...
	go build -o "$${GOPATH%%:*}/bin/msgp" .
	go generate ./msgp
	go generate ./_generated
	go test ./msgp
	go test ./_generated
