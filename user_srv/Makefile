.PHONY: build
build:
	go mod tidy
	go build server.go


.PHONY: clean
clean:
	rm server


.PHONY: proto
proto:
	cd proto && \
	protoc -I /Users/albelt/Gitlab/common/protos -I . \
		--go_out=. \
		--go-grpc_out=. \
		--validate_out="lang=go:." \
		const.proto msg.proto svc.proto


.PHONY: config
config:
	cd config && \
	protoc -I . --go_out=. config.proto && \
	protoc-go-inject-tag -input=config.pb.go