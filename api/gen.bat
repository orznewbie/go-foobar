
cd %~dp0%/../

protoc -I . -I ./third_party --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. ^
    api/test/data.proto api/test/user.proto
