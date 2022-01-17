# Testing out server

Start server:

go run server.go

Get 5 random values in range of 0 to 50:

grpc_cli call localhost:9000 GetValuePair "minRange: 0, maxRange: 50, count: 5"  --protofiles=protos/value.proto

Get random values every other second:

grpc_cli call localhost:9000 GetValueStream "minRange: 0, maxRange: 50"  --protofiles=protos/value.proto