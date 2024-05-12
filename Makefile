
install-dev:
	pre-commit install

dev-gateway:
	cd ./gateway && air

dev-orders:
	cd ./orders && air

proto-gen:
	protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    ./common/api/oms.proto
