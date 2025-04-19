GENERATOR = openapi-generator-cli
SPEC = openapi.yaml
OUT_DIR = generated

generate:
	$(GENERATOR) generate \
		-i $(SPEC) \
		-g go-server \
		-o $(OUT_DIR) \
		--additional-properties=packageName=schedules,serverPort=5252

clean:
	rm -rf $(OUT_DIR)