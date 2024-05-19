// main.go
package main

import (
    "context"
    "log"
    "net"
    "net/http"
    "google.golang.org/grpc"
    "google.golang.org/grpc/reflection"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    pb "example.com/grpc-go/proto"
)

var (
    grpcRequests = prometheus.NewCounter(prometheus.CounterOpts{
        Name: "masters_grpc_requests",
        Help: "Total number of gRPC requests received",
    })
)

type server struct {
    pb.UnimplementedMasterServiceServer
}

func (s *server) GetMasterData(ctx context.Context, req *pb.MasterRequest) (*pb.MasterResponse, error) {
    grpcRequests.Inc()
    return &pb.MasterResponse{Data: "Hello " + req.Query}, nil
}

func main() {
    prometheus.MustRegister(grpcRequests)
    
    lis, err := net.Listen("tcp", ":8082")
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }
    s := grpc.NewServer()
    pb.RegisterMasterServiceServer(s, &server{})

    // Register reflection service on gRPC server.
    reflection.Register(s)

    go func() {
        http.Handle("/metrics", promhttp.Handler())
        http.ListenAndServe(":8080", nil)
    }()

    if err := s.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}
