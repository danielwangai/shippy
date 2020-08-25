package shippy_service_consignment

import(
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"sync"

	pb "github.com/shippy-service-consignment/proto/consignment"
)

const (
	port = ":50051"
)

type repository interface {
	Create(*pb.Consignment) (*pb.Consignment, error)
}

type Repository struct {
	mu sync.Mutex
	consignments []*pb.Consignment
}

// Create a new consignment
func (repo *Repository) Create(consignment *pb.Consignment) (*pb.Consignment, error) {
	repo.mu.Lock()
	repo.consignments = append(repo.consignments, consignment)
	repo.mu.Unlock()
	return consignment, nil
}

type service struct {
	repo repository
}

func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment) (*pb.Response, error) {
	consignment, err := s.repo.Create(req)
	if err != nil {
		return nil, err
	}
	return &pb.Response{Created: true, Consignment: consignment}, nil
}

func main() {
	repo := &Repository{}
	// set up gRPC server
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	// Register service on the server
	pb.RegisterShippingServiceServer(s, &service{repo})
	// register reflection service on gRPC server.
	reflection.Register(s)
	log.Println("Running on port ", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
