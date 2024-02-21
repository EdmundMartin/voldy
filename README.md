# Voldy

An attempted reimplementation of LinkedIn's Voldermort data store.


### Generate Protobufs
```
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    pkg/protocol/voldy.proto
```

## TODO
* Fix usage of versioning - introduce client which supports - get prior to put operations
* Read-repair implementation
* Postgres backed storage engine
* Look at implementing full DynamoDB API - HashKey/Sort key - should be easy with B-tree storage engine.
