package main

import (
	"context"

	"github.com/bcc-code/mediabank-bridge/log"
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
	log.L.Debug().Msg("CreateSubclip")
	err := s.vantageClient.CreateSubclip(vantage.CreateSubclipParams{
		In:      req.In,
		Out:     req.Out,
		AssetID: req.AssetId,
		Title:   req.Title,
	})

	log.L.Debug().Err(err)

	return &proto.CreateSubclipResponse{
		Message: "X",
	}, nil
}
