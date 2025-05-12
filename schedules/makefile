GENERATOR = oapi-codegen
SPEC = openapi.yaml
OPENAPI_OUT_DIR = pkg/rest
PROTO_DIR = internal/server/generated/grpc/
PROTO_FILE = schedule.proto

generateOpenapi:
	$(GENERATOR) -generate types -package rest openapi.yaml > $(OPENAPI_OUT_DIR)/types.gen.go

generateProto:
	protoc --proto_path=$(PROTO_DIR) \
               --go_out=$(PROTO_DIR) \
               --go-grpc_out=$(PROTO_DIR) \
               $(PROTO_DIR)/schedule.proto

clean:
	rm -rf $(API_OUT_DIR)

test:
	go test .\...

