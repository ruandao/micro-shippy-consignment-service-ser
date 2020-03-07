package main

import (
	"context"
	"fmt"
	"github.com/micro/go-micro"
	"github.com/ruandao/micro-shippy-consignment-service-ser/consignmentMongo"
	pb "github.com/ruandao/micro-shippy-consignment-service-ser/proto/consignment"
	"github.com/ruandao/micro-shippy-vessel-service-ser/lib"
	vesselProto "github.com/ruandao/micro-shippy-vessel-service-ser/proto/vessel"
	"log"
	"os"
)

const (
	defaultHost = "datastore:27017"
)

func main() {
	
	// Create a new service. Optionally include some options here.
	srv := micro.NewService(
		micro.Name(consignmentMongo.CONST_SER_NAME_CONSIGNMENT),
	)

	// Init will parse the command line flags.
	srv.Init()
	
	uri := os.Getenv("DB_HOST")
	if uri == "" {
		uri = defaultHost
	}
	
	client, err := consignmentMongo.CreateClient(context.Background(), uri, 0)
	if err != nil {
		log.Panic(err)
	}
	defer client.Disconnect(context.Background())

	consignmentCollection := client.Database("shippy").Collection("consignments")

	repository := &consignmentMongo.MongoRepository{Collection: consignmentCollection}
	vesselClient := vesselProto.NewVesselServiceClient(lib.CONST_SERVICE_NAME, srv.Client())
	h := &consignmentMongo.Handler{Repository: repository, VesselClient: vesselClient}


	// Register handler
	pb.RegisterShippingServiceHandler(srv.Server(), h)

	// Run the server
	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}
