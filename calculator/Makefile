.PHONY: proto
proto:
	@protoc -I . calculatorpb/calculator.proto --go_out=plugins=grpc:.

.PHONY: server
server:
	@go run server/server.go

.PHONY: client
client:
	@go run client/client.go
