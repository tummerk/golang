GENERATOR = openapi-generator-cli
SPEC = openapi.yaml
API_OUT_DIR = generatedOpenapi
PROTO_DIR = ./generatedProtobuff/proto
PROTO_GEN_DIR = ./generatedProtobuff/gen
PROTO_FILE = schedule.proto

generateOpenapi:
	$(GENERATOR) generate \
		-i $(SPEC) \
		-g go-server \
		-o generatedOpenapi
		--additional-properties=packageName=schedules,serverPort=5252

generateProto:
	protoc --proto_path=./generatedProtobuff/proto \
               --go_out=./generatedProtobuff/gen \
               --go-grpc_out=./generatedProtobuff/gen \
               ./generatedProtobuff/proto/schedule.proto

clean:
	rm -rf $(API_OUT_DIR)