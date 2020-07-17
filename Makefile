CURRENT_DIR=$(shell pwd)
PROTO_DIR=${CURRENT_DIR}/proto
PROTO_GEN_DIR=${PROTO_DIR}/plugin

grpc-build:
	@echo ">> Building proto plugin..."
	@protoc -I=${PROTO_DIR} ${PROTO_DIR}/plugin.proto --go_out=plugins=grpc:${PROTO_GEN_DIR}
	@echo ">> Building proto plugin done"
	@echo ">> --------------------------"

grpc-mock:
	@echo ">> Mocking grpc..."
	@mockery -dir=${PROTO_GEN_DIR} -name=PluginClient -output=${PROTO_GEN_DIR}/mocks -quiet
	@echo ">> Mocking done"
	@echo ">> --------------------------"

gomod-tidy:
	@echo ">> Tidy go.mod..."
	@go mod tidy
	@echo ">> Tidy go.mod done"
	@echo ">> --------------------------"

build: gomod-tidy grpc-build grpc-mock