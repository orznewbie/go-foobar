@echo off
cd %~dp0%/../

protoc -I ./third_party -I %GOPATH%/src --go_out=paths=source_relative:%GOPATH%/src ^
	go-foobar/api/user/user.proto

protoc -I ./third_party -I %GOPATH%/src --go-grpc_out %GOPATH%/src --go-grpc_opt paths=source_relative ^
	go-foobar/api/user/user.proto

protoc -I ./third_party -I %GOPATH%/src --grpc-gateway_out %GOPATH%/src ^
	--grpc-gateway_opt logtostderr=true ^
	--grpc-gateway_opt paths=source_relative ^
	go-foobar/api/user/user.proto

protoc -I ./third_party -I %GOPATH%/src --validate_out=paths=source_relative,lang=go:%GOPATH%/src ^
	go-foobar/api/user/user.proto

protoc -I ./third_party -I %GOPATH%/src --openapiv2_out %GOPATH%/src ^
	--openapiv2_opt logtostderr=true ^
	--openapiv2_opt json_names_for_fields=true ^
	--openapiv2_opt openapi_naming_strategy=simple ^
	go-foobar/api/user/user.proto

cd %~dp0%