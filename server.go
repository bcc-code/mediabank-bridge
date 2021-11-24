package main

import (
	"context"

	"github.com/bcc-code/mediabank-bridge/proto"
	"github.com/bcc-code/mediabank-bridge/vantage"

	"github.com/davecgh/go-spew/spew"
)

// Server implements all GRPC calls
type Server struct {
	proto.UnimplementedMediabankBridgeServer

	vantageClient *vantage.Client
}

// CreateSubclip from the specified ID
func (s Server) CreateSubclip(ctx context.Context, req *proto.CreateSubclipRequest) (*proto.CreateSubclipResponse, error) {
	spew.Dump(req)
	err := s.vantageClient.CreateSubclip()
	return &proto.CreateSubclipResponse{}, err
}
