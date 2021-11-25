package main

import (
	"context"

	"github.com/bcc-code/mediabank-bridge/proto"
	"github.com/bcc-code/mediabank-bridge/vantage"
)

// Server implements all GRPC calls
type Server struct {
	proto.UnimplementedMediabankBridgeServer

	vantageClient *vantage.Client
}

// CreateSubclip from the specified ID
func (s Server) CreateSubclip(ctx context.Context, req *proto.CreateSubclipRequest) (*proto.CreateSubclipResponse, error) {
	err := s.vantageClient.CreateSubclip(vantage.CreateSubclipParams{
		In:      req.In,
		Out:     req.Out,
		AssetID: req.AssetId,
		Title:   req.Title,
	})

	return &proto.CreateSubclipResponse{}, err
}
