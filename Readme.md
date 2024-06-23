## A Voting system


To generate proto files
```shell
protoc -I ./proto \
  --go_out ./proto --go_opt paths=source_relative \
  --go-grpc_out ./proto --go-grpc_opt paths=source_relative \
  --grpc-gateway_out ./proto --grpc-gateway_opt paths=source_relative \
--plugin=protoc-gen-grpc-gateway=${GOPATH}/bin/protoc-gen-grpc-gateway \
  ./proto/user.proto
```