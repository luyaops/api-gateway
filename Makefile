register:
    protoc -I$GOPATH/src -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis $GOPATH/src/github.com/luyaops/*/proto/*.proto --parse_out=.