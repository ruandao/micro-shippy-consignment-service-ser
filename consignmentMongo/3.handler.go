package consignmentMongo

import (
	"context"
	"errors"
	pb "github.com/ruandao/micro-shippy-consignment-service-ser/proto/consignment"
	vesselProto "github.com/ruandao/micro-shippy-vessel-service-ser/proto/vessel"
	"log"
)

type Handler struct {
	Repository
	VesselClient vesselProto.VesselServiceClient
}

// CreateConsignment - we created just one method on our service,
// which is a create method, which takes a context and a request as an
// argument, these are handled by the gRPC server.
func (s *Handler) CreateConsignment(ctx context.Context, req *pb.Consignment, res *pb.Response) error {

	// Here we call a client instance of our vessel service with our consignment weight,
	// and the amount of containers as the capacity value
	log.Printf("run here 1")
	vesselResponse, err := s.VesselClient.FindAvailable(ctx, &vesselProto.Specification{
		Capacity:  req.Weight,
		MaxWeight: int32(len(req.Containers)),
	})
	log.Printf("run here 2")
	if vesselResponse == nil {
		return errors.New("error fetching vessel, returned nil")
	}

	log.Printf("run here 3")
	if err != nil {
		return err
	}

	log.Printf("run here 4")
	// We set the VesselId as the vessel we got back from our
	// vessel service
	req.VesselId = vesselResponse.Vessel.Id

	// Save our consignment
	if err = s.Repository.Create(ctx, MarshalConsignment(req)); err != nil {
		return err
	}
	log.Printf("run here 5")
	res.Created = true
	res.Consignment = req
	log.Printf("run here 6")
	return nil
}

// getConsignments -
func (s *Handler) GetConsignments(ctx context.Context, req *pb.GetRequest, res *pb.Response) error {
	consignments, err := s.Repository.GetAll(ctx)
	if err != nil {
		return err
	}

	res.Consignments = UnmarshalConsignmentCollection(consignments)
	return nil
}