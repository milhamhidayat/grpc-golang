.PHONY: proto
proto:
	@protoc -I . pb/pb.proto --go_out=plugins=grpc:.

.PHONY: client
client:
	@go run client/main.go

.PHONY: server
server:
	@go run server/main.go
