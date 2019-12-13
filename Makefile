.PHONY: build run

proto/consignment/consignment.pb.go: proto
	protoc -I. --go_out=plugins=micro:. proto/consignment/consignment.proto

consignment-service: main.go proto
	GOOS=linux GOARCH=amd64 go build

.build/.docker-container.stamp: Dockerfile main.go proto/consignment/consignment.pb.go go.mod go.sum
	docker build -t shippy-service-consignment .
	mkdir -p $(dir $@)
	touch $@

build: .build/.docker-container.stamp

run: build
	docker run -p 50051:50051 \
		-e MICRO_SERVER_ADDRESS=:50051 \
		shippy-service-consignment
