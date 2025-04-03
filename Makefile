.PHONY: build-all clean swagger migrate protoc-bins
PROTOC_FILES = $(shell find . -type f -name '*.proto')
PB_FILES = $(patsubst %.proto, %.pb.go, $(PROTOC_FILES))
PB_STORM_FILES = $(patsubst %.proto, %.pb.storm.go, $(PROTOC_FILES))
BIN = $(addprefix bin/,$(shell ls cmd))

BININFO_VERSION     ?= $(shell git describe --tags --always --abbrev=7)
BININFO_COMMIT_HASH ?= $(shell git rev-parse --verify HEAD)
BININFO_BUILD_DATE  ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

GOPATH = $(PWD)/gopath
PROTOC_BINS = $(shell find $(PWD)/bin -type f -name 'protoc*')

MIGRATE = $(PWD)/bin/migrate
PROTOC  = $(PWD)/bin/protoc
SWAGGER = $(PWD)/bin/swagger

MIGRATE_VERSION             := v4.18.1
PROTOC_VERSION              := 30.2
PROTOC_GEN_GO_VERSION       := v1.36.6
PROTOC_GEN_GO_GRPC_VERSION  := cdbdb759dd67c89544f9081f854c284493b5461c # v1.71.1
PROTOC_GEN_STORMRPC_VERSION := v0.4.2
SWAGGER_VERSION             := v0.30.4

MIGRATE_PKG_URL            := "https://github.com/golang-migrate/migrate/releases/download/$(MIGRATE_VERSION)/migrate.darwin-amd64.tar.gz"
PROTOC_GEN_GO_PKG_URL      := "https://github.com/protocolbuffers/protobuf-go/releases/download/$(PROTOC_GEN_GO_VERSION)/protoc-gen-go.$(PROTOC_GEN_GO_VERSION).darwin.arm64.tar.gz"

TMP_DIR := $(PWD)/tmp

build-all: $(PROTOC_BINS) $(PB_FILES) $(PB_STORM_FILES) $(BIN)

LDFLAGS= -ldflags="-X 'github.com/sapcc/go-api-declarations/bininfo.buildDate=$(BININFO_BUILD_DATE)' -X 'github.com/sapcc/go-api-declarations/bininfo.commit=$(BININFO_COMMIT_HASH)' -X 'github.com/sapcc/go-api-declarations/bininfo.version=$(BININFO_VERSION)'"

tmp:
	mkdir tmp

tmp/protobuf.zip: tmp
	curl -L -o $(TMP_DIR)/protobuf.zip https://github.com/protocolbuffers/protobuf/releases/download/v$(PROTOC_VERSION)/protoc-$(PROTOC_VERSION)-osx-aarch_64.zip

bin/protoc: tmp/protobuf.zip
	mkdir $(TMP_DIR)/protobuf-unzipped
	unzip $(TMP_DIR)/protobuf.zip -d $(TMP_DIR)/protobuf-unzipped
	cp $(TMP_DIR)/protobuf-unzipped/bin/protoc $(PWD)/bin
	rm -rf $(TMP_DIR)/protobuf-unzipped

bin/migrate:
	$(PWD)/.bin/install_tarball_pkg.sh $(MIGRATE_PKG_URL) $(TMP_DIR) $(PWD)/bin/migrate

bin/protoc-gen-go:
	$(PWD)/.bin/install_tarball_pkg.sh $(PROTOC_GEN_GO_PKG_URL) $(TMP_DIR) $(PWD)/bin/protoc-gen-go

bin/protoc-gen-go-grpc: gopath
	GOPATH=$(GOPATH) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@$(PROTOC_GEN_GO_GRPC_VERSION)

bin/protoc-gen-stormrpc: gopath
	GOPATH=$(GOPATH) go install github.com/actatum/stormrpc/cmd/protoc-gen-stormrpc@$(PROTOC_GEN_STORMRPC_VERSION)

bin/swagger: gopath
	GOPATH=$(GOPATH) go install github.com/go-swagger/go-swagger/cmd/swagger@$(SWAGGER_VERSION)

bin/%: cmd/%/main.go
	go build $(LDFLAGS) -o $@ $<

swagger: gopath/bin/swagger
	$(SWAGGER) generate server --exclude-main --copyright-file COPYRIGHT.txt
	$(SWAGGER) generate model --copyright-file COPYRIGHT.txt
	$(SWAGGER) generate client --copyright-file COPYRIGHT.txt

protoc-bins: bin/protoc \
	gopath/bin/protoc-gen-go \
	gopath/bin/protoc-gen-go-grpc \
	gopath/bin/protoc-gen-stormrpc

markdown: gopath/bin/swagger
	$(SWAGGER) generate markdown --copyright-file COPYRIGHT.txt --output= docs/api.md

migrate: gopath/bin/migrate
	$(MIGRATE) -path db/migrations -database "cockroachdb://root@localhost:26257/andromeda?sslmode=disable" drop -f
	$(MIGRATE) -path db/migrations -database "cockroachdb://root@localhost:26257/andromeda?sslmode=disable" up

%.pb.storm.go: %.proto
	$(PROTOC) --stormrpc_out=. --stormrpc_opt=paths=source_relative -I. $<

%.pb.go: %.proto
	$(PROTOC) --go_out=. --go_opt=paths=source_relative -I. $<


clean:
	rm -f bin/*
	rm -rf $(TMP_DIR)/*
