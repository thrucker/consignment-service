ARG buildImage="golang:alpine"
FROM ${buildImage} as builder

RUN apk --no-cache add git protobuf
RUN go get -u github.com/micro/protobuf/proto
RUN go get -u github.com/micro/protobuf/protoc-gen-go

WORKDIR /app/shippy-service-consignment

COPY go.mod ./go.mod
COPY go.sum ./go.sum

RUN go mod download

COPY . .

RUN go generate
RUN CGO_ENABLED=0 GOOS=linux go build -installsuffix cgo -o shippy-service-consignment main.go datastore.go handler.go repository.go

FROM alpine:latest as main

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/shippy-service-consignment/shippy-service-consignment .

CMD ["./shippy-service-consignment"]

FROM builder as obj-cache

COPY --from=builder /root/.cache /root/.cache
