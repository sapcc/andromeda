.PHONY: build-all clean swagger migrate
PROTOC_FILES = $(shell find . -type f -name '*.proto')
PB_FILES = $(patsubst %.proto, %.pb.go, $(PROTOC_FILES))
PB_STORM_FILES = $(patsubst %.proto, %.pb.storm.go, $(PROTOC_FILES))
BIN = $(addprefix bin/,$(shell ls cmd))

BININFO_VERSION     ?= $(shell git describe --tags --always --abbrev=7)
BININFO_COMMIT_HASH ?= $(shell git rev-parse --verify HEAD)
BININFO_BUILD_DATE  ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

build-all: $(PB_FILES) $(PB_STORM_FILES) $(BIN)

LDFLAGS= -ldflags="-X 'github.com/sapcc/go-api-declarations/bininfo.buildDate=$(BININFO_BUILD_DATE)' -X 'github.com/sapcc/go-api-declarations/bininfo.commit=$(BININFO_COMMIT_HASH)' -X 'github.com/sapcc/go-api-declarations/bininfo.version=$(BININFO_VERSION)'"

bin/%: cmd/%/main.go
	go build $(LDFLAGS) -o $@ $<

swagger:
	swagger generate server --exclude-main --copyright-file COPYRIGHT.txt
	swagger generate model --copyright-file COPYRIGHT.txt
	swagger generate client --copyright-file COPYRIGHT.txt

markdown:
	swagger generate markdown --copyright-file COPYRIGHT.txt --output= docs/api.md

migrate:
	migrate -path db/migrations -database "cockroachdb://root@localhost:26257/andromeda?sslmode=disable" drop -f
	migrate -path db/migrations -database "cockroachdb://root@localhost:26257/andromeda?sslmode=disable" up

%.pb.storm.go: %.proto
	protoc --stormrpc_out=. --stormrpc_opt=paths=source_relative -I. $<

%.pb.go: %.proto
	protoc --go_out=. --go_opt=paths=source_relative -I. $<


clean:
	rm -f bin/*
