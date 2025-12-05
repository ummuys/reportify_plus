PROTO_DIR = proto
PB_OUT_DIR = pb

.PHONY: proto

proto:
	protoc -I $(PROTO_DIR) \
	--go_out=$(PB_OUT_DIR) --go_opt=paths=source_relative \
	--go-grpc_out=$(PB_OUT_DIR) --go-grpc_opt=paths=source_relative \
	$(PROTO_DIR)/auth/v1/*.proto