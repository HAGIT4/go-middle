generate-swagger:
	swag init --o ./docs -d ./internal/server/api -g router.go --parseDependency

build-proto:
	mkdir -p pb
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative --plugin=protoc-gen-go_grpc=${GOPATH}/bin/protoc-gen-go-grpc pb/metric.proto