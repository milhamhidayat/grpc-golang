.PHONY: proto
proto:
	@protoc -I . greetpb/greet.proto --go_out=plugins=grpc:.

.PHONY: client
client:
	@go run greet_client/main.go

.PHONY: server
server:
	@go run greet_server/main.go

