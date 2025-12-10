PROTO_DIR = proto
PB_OUT_DIR = pb

.PHONY: proto

proto-auth:
	protoc -I $(PROTO_DIR) \
	--go_out=$(PB_OUT_DIR) --go_opt=paths=source_relative \
	--go-grpc_out=$(PB_OUT_DIR) --go-grpc_opt=paths=source_relative \
	$(PROTO_DIR)/auth/v1/*.proto 

proto-report:
	protoc -I $(PROTO_DIR) \
	--go_out=$(PB_OUT_DIR) --go_opt=paths=source_relative \
	--go-grpc_out=$(PB_OUT_DIR) --go-grpc_opt=paths=source_relative \
	$(PROTO_DIR)/report/v1/*.proto

proto-all: proto-auth proto-report