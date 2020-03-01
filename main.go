package main

import (
	"context"
	"fmt"
	"github.com/micro/go-micro"
	pb "github.com/ruandao/micro-shippy-consignment-service-ser/proto/consignment"
	vesselProto "github.com/ruandao/micro-shippy-vessel-service/proto/vessel"
	"log"
	"os"
)

const (
	defaultHost = "datastore:27017"
)

func main() {
	
	// Create a new service. Optionally include some options here.
	srv := micro.NewService(
		micro.Name("go.micro.srv.consignment"),
	)

	// Init will parse the command line flags.
	srv.Init()
	
	uri := os.Getenv("DB_HOST")
	if uri == "" {
		uri = defaultHost
	}
	
	client, err := CreateClient(context.Background(), uri, 0)
	if err != nil {
		log.Panic(err)
	}
	defer client.Disconnect(context.Background())

	consignmentCollection := client.Database("shippy").Collection("consignments")

	repository := &MongoRepository{consignmentCollection}
	vesselClient := vesselProto.NewVesselServiceClient("go.micro.srv.vessel", srv.Client())
	h := &handler{repository, vesselClient}


	// Register handler
	pb.RegisterShippingServiceHandler(srv.Server(), h)

	// Run the server
	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}
