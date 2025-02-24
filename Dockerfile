FROM golang:1.24-alpine AS builder
RUN apk add --no-cache make gcc musl-dev protoc git && go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28 && go install github.com/actatum/stormrpc/cmd/protoc-gen-stormrpc@latest

COPY . /src
RUN make -C /src

################################################################################

FROM alpine:3.21
LABEL source_repository="https://github.com/sapcc/andromeda"

RUN apk add --no-cache ca-certificates
COPY . /src
COPY --from=builder /src/bin/ /usr/bin/
ENTRYPOINT [ "/usr/bin/andromeda-server" ]
