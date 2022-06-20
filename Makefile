.PHONY: build-all clean swagger migrate
PROTOC_FILES = $(shell find . -type f -name '*.proto')
PB_FILES = $(patsubst %.proto, %.pb.go, $(PROTOC_FILES))
PB_MICRO_FILES = $(patsubst %.proto, %.pb.micro.go, $(PROTOC_FILES))

build-all: $(PB_FILES) $(PB_MICRO_FILES) bin/andromeda bin/m31ctl

bin/andromeda: main.go $(PB_FILES) $(PB_MICRO_FILES)
	go build -o bin/andromeda

bin/m31ctl: cmd/m31ctl.go internal/cli/*.go
	go build -o bin/m31ctl cmd/m31ctl.go

swagger:
	swagger generate server --exclude-main --copyright-file COPYRIGHT.txt
	swagger generate model --copyright-file COPYRIGHT.txt
	swagger generate client --copyright-file COPYRIGHT.txt

markdown:
	swagger generate markdown --copyright-file COPYRIGHT.txt --output= docs/api.md

migrate:
	migrate -path db/migrations -database "cockroachdb://root@localhost:26257/andromeda?sslmode=disable" drop -f
	migrate -path db/migrations -database "cockroachdb://root@localhost:26257/andromeda?sslmode=disable" up

%.pb.micro.go: %.proto
	protoc --micro_out=../../.. -I. $<

%.pb.go: %.proto
	protoc --go_out=../../.. -I. $<


clean:
	rm -f bin/*
