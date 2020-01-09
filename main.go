package main

import (
	"context"
	"fmt"
	"github.com/micro/go-micro"
	pb "github.com/thrucker/consignment-service/proto/consignment"
	"github.com/thrucker/vessel-service/proto/vessel"
	"log"
	"os"
)

const (
	port = ":50051"
	defaultHost = "datastore:27017"
)

func main() {
	srv := micro.NewService(
		micro.Name("shippy.service.consignment"),
	)

	srv.Init()

	uri := os.Getenv("DB_HOST")
	if uri == "" {
		uri = defaultHost
	}

	client, err := CreateClient(uri)
	if err != nil {
		log.Panic(err)
	}
	defer client.Disconnect(context.Background())

	consignmentCollection := client.Database("shippy").Collection("consignments")

	repository := &MongoRepository{consignmentCollection}
	vesselClient := vessel.NewVesselServiceClient("shippy.service.vessel", srv.Client())
	h := &handler{repository, vesselClient}

	pb.RegisterShippingServiceHandler(srv.Server(), h)

	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}
