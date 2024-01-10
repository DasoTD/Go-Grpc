package api

import (

	db "github.com/dasotd/go_grpc/db/sqlc"
	"github.com/dasotd/go_grpc/pb"
)


type Server struct {
	pb.UnimplementedAccountAPIServer
	grpcDB           db.grpc
}

func NewServer( grpcDB db.grpc) (*Server, error) {
	
	server := &Server{
		grpcDB:           grpcDB,
	}

	return server, nil
}
