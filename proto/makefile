PROTO_FILES := $(shell find . shared-proto -name '*.proto')

.PHONY: proto clean

proto:
	mkdir -p generated
	protoc \
        -I . \
        --go_out=generated \
        --go_opt=paths=source_relative \
        --go-grpc_out=generated \
        --go-grpc_opt=paths=source_relative \
		$(PROTO_FILES)
	
	cd generated && \
		go mod init proto && \
		go mod tidy

clean:
	rm -rf generated