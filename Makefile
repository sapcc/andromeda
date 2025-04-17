.PHONY: build-all clean swagger migrate build-akamai-dns-requests
PROTOC_FILES = $(shell find . -type f -name '*.proto')
PB_FILES = $(patsubst %.proto, %.pb.go, $(PROTOC_FILES))
PB_STORM_FILES = $(patsubst %.proto, %.pb.storm.go, $(PROTOC_FILES))
BIN = $(addprefix bin/,$(shell ls cmd))
BIN_DEPS = $(addprefix bin/,$(shell echo migrate protoc protoc-gen-go protoc-gen-go-grpc protoc-gen-stormrpc swagger))

BININFO_VERSION     ?= $(shell git describe --tags --always --abbrev=7)
BININFO_COMMIT_HASH ?= $(shell git rev-parse --verify HEAD)
BININFO_BUILD_DATE  ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

MIGRATE ?= $(PWD)/bin/migrate
PROTOC  ?= $(PWD)/bin/protoc
SWAGGER ?= $(PWD)/bin/swagger

MIGRATE_VERSION             ?= v4.18.1
PROTOC_VERSION              ?= 30.2
PROTOC_GEN_GO_VERSION       ?= v1.36.6
PROTOC_GEN_GO_GRPC_VERSION  ?= cdbdb759dd67c89544f9081f854c284493b5461c # v1.71.1
PROTOC_GEN_STORMRPC_VERSION ?= v0.4.2
SWAGGER_VERSION             ?= v0.30.4

TMP_DIR := $(PWD)/tmp

build-all: $(BIN_DEPS) $(PB_FILES) $(PB_STORM_FILES) $(BIN)
ci-build-all: $(PB_FILES) $(PB_STORM_FILES) $(BIN)

# Specific target for building the Akamai DNS requests CLI tool
build-akamai-dns-requests:
	go build $(LDFLAGS) -o build/andromeda-akamai-total-dns-requests cmd/andromeda-akamai-total-dns-requests/main.go

LDFLAGS= -ldflags="-X 'github.com/sapcc/go-api-declarations/bininfo.buildDate=$(BININFO_BUILD_DATE)' -X 'github.com/sapcc/go-api-declarations/bininfo.commit=$(BININFO_COMMIT_HASH)' -X 'github.com/sapcc/go-api-declarations/bininfo.version=$(BININFO_VERSION)'"

bin/migrate:
	go install github.com/golang-migrate/migrate/v4/cmd/migrate@$(MIGRATE_VERSION)
	-cp $(GOPATH)/bin/migrate $(PWD)/bin

bin/protoc-gen-go:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@$(PROTOC_GEN_GO_VERSION) && cp $(GOPATH)/bin/protoc-gen-go $(PWD)/bin

bin/protoc-gen-go-grpc:
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@$(PROTOC_GEN_GO_GRPC_VERSION) && cp $(GOPATH)/bin/protoc-gen-go-grpc $(PWD)/bin

bin/protoc-gen-stormrpc:
	go install github.com/actatum/stormrpc/cmd/protoc-gen-stormrpc@$(PROTOC_GEN_STORMRPC_VERSION) && cp $(GOPATH)/bin/protoc-gen-stormrpc $(PWD)/bin

bin/swagger:
	go install github.com/go-swagger/go-swagger/cmd/swagger@$(SWAGGER_VERSION) && cp $(GOPATH)/bin/swagger $(PWD)/bin

tmp:
	test -d $(TMP_DIR) || mkdir $(TMP_DIR)

tmp/protobuf.zip: tmp
	curl -L -o $(TMP_DIR)/protobuf.zip https://github.com/protocolbuffers/protobuf/releases/download/v$(PROTOC_VERSION)/protoc-$(PROTOC_VERSION)-osx-aarch_64.zip

bin/protoc: tmp/protobuf.zip
	mkdir $(TMP_DIR)/protobuf-unzipped
	unzip $(TMP_DIR)/protobuf.zip -d $(TMP_DIR)/protobuf-unzipped
	cp $(TMP_DIR)/protobuf-unzipped/bin/protoc $(PWD)/bin
	chmod 755 $(PWD)/bin/protoc
	rm -rf $(TMP_DIR)/protobuf-unzipped

bin/%: cmd/%/main.go
	go build $(LDFLAGS) -o $@ $<

swagger: bin/swagger
	$(SWAGGER) generate server --exclude-main --copyright-file COPYRIGHT.txt
	$(SWAGGER) generate client --copyright-file COPYRIGHT.txt
	$(SWAGGER) generate model --copyright-file COPYRIGHT.txt --struct-tags db

markdown: bin/swagger
	$(SWAGGER) generate markdown --copyright-file COPYRIGHT.txt --output docs/api.md

serve-swagger-docs: bin/swagger
	$(SWAGGER) serve swagger.yml -p 9900

migrate: bin/migrate
	$(MIGRATE) -path db/migrations -database "cockroachdb://root@localhost:26257/andromeda?sslmode=disable" drop -f
	$(MIGRATE) -path db/migrations -database "cockroachdb://root@localhost:26257/andromeda?sslmode=disable" up

%.pb.storm.go: %.proto
	$(PROTOC) --stormrpc_out=. --stormrpc_opt=paths=source_relative -I. $<

%.pb.go: %.proto
	$(PROTOC) --go_out=. --go_opt=paths=source_relative -I. $<


clean:
	rm -f bin/*
	rm -rf $(TMP_DIR)
