package agent

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/srl-wim/srl-ndk-git/api/proto/gitapi"
)

// StartGRPCServer start a gRPC server and waits for connection
func (a *Agent) StartGRPCServer() {
	// create a listener on TCP port 7777
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 7777))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// create a server instance
	s := a.GRPCServer{}
	// create a gRPC server object
	grpcServer := grpc.NewServer()
	// attach the Ping service to the server
	gitapi.RegisterGitServer(grpcServer, &s)
	// start the server
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

// GRPCServer executes the command and generates a response
func (a *Agent) GRPCServer(ctx context.Context, in *gitapi.Action) (*gitapi.ActionResponse, error) {
	log.Printf("Receive message %s, %s", in.Kind, in.Attributes)
	return &ActionResponse{Response: fmt.Sprintf("Message %s processed ok", in.Kind)}, nil
}
