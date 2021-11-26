package main

//asldkjalskdlkasjd
import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/bcc-code/mediabank-bridge/auth0"
	"github.com/bcc-code/mediabank-bridge/config"
	"github.com/bcc-code/mediabank-bridge/log"
	"github.com/bcc-code/mediabank-bridge/proto"
	"github.com/bcc-code/mediabank-bridge/vantage"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

// Authorization unary interceptor function to handle authorize per RPC call
func serverInterceptor(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "Failed to get metadata")
	}

	authHeader, ok := md["token"]
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "missing token")
	}

	token := authHeader[0]

	err := auth0.ValidateJWT(token, auth0.JWTConfig{
		Audience: "media.bcc.mediabanken",
		Issuer:   "https://login.bcc.no/",
		Domain:   "login.bcc.no",
	})

	if err != nil {
		log.L.Info().Err(err).Msg("Auth failed")
		return nil, status.Errorf(codes.Unauthenticated, err.Error())
	}

	return handler(ctx, req)
}

func main() {
	log.ConfigureGlobalLogger(zerolog.DebugLevel)

	port := os.Getenv("PORT")

	if port == "" {
		port = "8083"
	}
	configFile := os.Getenv("CONF_FILE")
	conf := config.MustReadConfigFile(configFile)

	vantageClient, err := vantage.NewClient(conf.Vantage)

	if err != nil {
		panic(err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.L.Fatal().
			Str("PORT", port).
			Err(err).
			Msgf("failed to listen")
	}

	s := grpc.NewServer(grpc.UnaryInterceptor(serverInterceptor))
	reflection.Register(s) // This is needed for using most debug UIs

	proto.RegisterMediabankBridgeServer(s, &Server{
		vantageClient: vantageClient,
	})

	log.L.Info().Msgf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.L.Fatal().
			Str("PORT", port).
			Err(err).
			Msgf("failed to serve")
	}
}
