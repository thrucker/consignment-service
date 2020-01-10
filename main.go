//go:generate protoc -I. --go_out=plugins=micro:. proto/consignment/consignment.proto
package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/metadata"
	"github.com/micro/go-micro/server"
	"github.com/thrucker/consignment-service/id"
	pb "github.com/thrucker/consignment-service/proto/consignment"
	userServiceId "github.com/thrucker/user-service/id"
	userService "github.com/thrucker/user-service/proto/user"
	vesselServiceId "github.com/thrucker/vessel-service/id"
	vesselService "github.com/thrucker/vessel-service/proto/vessel"
	"log"
	"os"
)

const (
	port        = ":50051"
	defaultHost = "datastore:27017"
)

func AuthWrapper(fn server.HandlerFunc) server.HandlerFunc {
	authClient := userService.NewUserServiceClient(userServiceId.UserServiceId, client.DefaultClient)

	return func(ctx context.Context, req server.Request, resp interface{}) error {
		meta, ok := metadata.FromContext(ctx)
		if !ok {
			return errors.New("no auth meta-data found in request")
		}

		token := meta["Token"]
		log.Println("Authenticating with token: ", token)

		_, err := authClient.ValidateToken(context.Background(), &userService.Token{Token: token})
		if err != nil {
			return err
		}
		log.Println("Token is valid")
		err = fn(ctx, req, resp)
		return err
	}
}

func main() {
	srv := micro.NewService(
		micro.Name(id.ConsignmentServiceId),
		micro.Version("latest"),
		micro.WrapHandler(AuthWrapper),
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
	vesselClient := vesselService.NewVesselServiceClient(vesselServiceId.VesselServiceId, srv.Client())
	h := &handler{repository, vesselClient}

	pb.RegisterShippingServiceHandler(srv.Server(), h)

	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}
