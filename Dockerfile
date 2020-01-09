FROM golang:alpine as builder

WORKDIR /app/shippy-service-consignment

COPY go.mod ./go.mod
COPY go.sum ./go.sum

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o shippy-service-consignment main.go datastore.go handler.go repository.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/shippy-service-consignment/shippy-service-consignment .

CMD ["./shippy-service-consignment"]
