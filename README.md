# Voldy

An attempted reimplementation of LinkedIn's Voldermort data store.


### Generate Protobufs
```
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    pkg/protocol/voldy.proto
```