.PHONY: proto
proto:
	@protoc -I . pb/blog.proto --go_out=plugins=grpc:.

.PHONY: client
client:
	@go run client/main.go

.PHONY: server
server:
	@go run server/main.go

.PHONY: mongo-up
mongo-up:
	@docker-compose up -d mongo

.PHONY: mongo-down
mongo-down:
	@docker stop avalon_mongo
