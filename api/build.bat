cd %~dp0%/../
@echo off

protoc -I . -I ./third_party ^
    --go_out paths=source_relative:. ^
    --go-grpc_out paths=source_relative:. ^
    --grpc-gateway_out . ^
    --grpc-gateway_opt logtostderr=true ^
    --grpc-gateway_opt paths=source_relative ^
    api/user/v1/user.proto

protoc -I . -I ./third_party ^
    --grpc-gateway_out ./api ^
    --grpc-gateway_opt logtostderr=true ^
    --grpc-gateway_opt paths=source_relative ^
    --grpc-gateway_opt grpc_api_configuration=api/google/longrunning/gateway_config.yaml ^
    --grpc-gateway_opt standalone=true ^
    google/longrunning/operations.proto