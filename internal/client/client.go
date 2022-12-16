// Package client contains grpc client and interceptors.
package client

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/vstebletsov89/go-developer-course-gophkeeper/internal/config"
	pb "github.com/vstebletsov89/go-developer-course-gophkeeper/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// RunClient starts client application to communicate with the user.
func RunClient(cfg *config.Config) error {
	// init global logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// default log level is info
	zerolog.SetGlobalLevel(config.ParseLogLevel(cfg.LogLevel))

	// debug config
	log.Debug().Msgf("%+v\n\n", cfg)

	// TODO: implement cli

	//ctx, cancel := context.WithCancel(context.Background())
	//defer cancel()

	//authClient, secretClient, err := startClient(cfg)
	//if err != nil {
	//	return err
	//}

	return nil
}

func startClient(cfg *config.Config) (*AuthClient, *SecretClient, error) {
	var clientConn *grpc.ClientConn
	authClient := NewAuthClient()
	secretClient := NewSecretClient()

	// start GRPC client with/without TLS
	var err error
	if cfg.EnableTLS {
		// Client using TLS credentials
		log.Info().Msg("GRPC client configuration with TLS credentials")
		transportCredentials, err := credentials.NewClientTLSFromFile("cert.pem", "")
		if err != nil {
			log.Error().Msgf("GRPC client credentials.NewClientTLSFromFile: %v", err.Error())
			return nil, nil, err
		}

		clientConn, err = grpc.Dial(cfg.ServerAddress, grpc.WithTransportCredentials(transportCredentials),
			grpc.WithUnaryInterceptor(authClient.UnaryInterceptorClient))
		if err != nil {
			log.Error().Msgf("GRPC client Dial: %v", err.Error())
			return nil, nil, err
		}
	} else {
		// client without TLS credentials
		log.Info().Msg("GRPC client configuration without TLS credentials")
		clientConn, err = grpc.Dial(cfg.ServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithUnaryInterceptor(authClient.UnaryInterceptorClient))
		if err != nil {
			log.Error().Msgf("GRPC client Dial: %v", err.Error())
			return nil, nil, err
		}
	}

	authClient.SetService(pb.NewAuthClient(clientConn))
	secretClient.SetService(pb.NewGophkeeperClient(clientConn))
	return authClient, secretClient, nil
}
