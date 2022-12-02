FROM golang:1.19-alpine as builder
RUN apk add --no-cache make gcc musl-dev protoc git && go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28 && go install github.com/go-micro/generator/cmd/protoc-gen-micro@latest

COPY . /src
RUN make -C /src

################################################################################

FROM alpine:3.17
LABEL source_repository="https://github.com/sapcc/andromeda"

RUN apk add --no-cache ca-certificates
COPY --from=builder /src/bin/ /usr/bin/
ENTRYPOINT [ "/usr/bin/andromeda-server" ]
