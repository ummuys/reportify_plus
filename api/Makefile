PROTO_DIR = proto
PB_OUT_DIR = pb

.PHONY: proto-auth proto-report proto-datasource proto-all ensure-pb-dirs

ensure-pb-dirs:
	@mkdir -p $(PB_OUT_DIR)

proto-auth: ensure-pb-dirs
	protoc -I $(PROTO_DIR) \
	--go_out=$(PB_OUT_DIR) --go_opt=paths=source_relative \
	--go-grpc_out=$(PB_OUT_DIR) --go-grpc_opt=paths=source_relative \
	$(PROTO_DIR)/auth/service/v1/*.proto

proto-report: ensure-pb-dirs
	protoc -I $(PROTO_DIR) \
	--go_out=$(PB_OUT_DIR) --go_opt=paths=source_relative \
	--go-grpc_out=$(PB_OUT_DIR) --go-grpc_opt=paths=source_relative \
	$(PROTO_DIR)/report/service/v1/*.proto

proto-datasource: ensure-pb-dirs
	protoc -I $(PROTO_DIR) \
	--go_out=$(PB_OUT_DIR) --go_opt=paths=source_relative \
	--go-grpc_out=$(PB_OUT_DIR) --go-grpc_opt=paths=source_relative \
	$(PROTO_DIR)/datasource/service/v1/*.proto

proto-all: proto-auth proto-report proto-datasource
