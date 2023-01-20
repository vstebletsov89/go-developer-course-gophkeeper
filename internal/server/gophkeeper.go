// Package server contains implementation of grpc server and interceptors.
package server

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq" // load postgres driver
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/vstebletsov89/go-developer-course-gophkeeper/internal/secure"
	"github.com/vstebletsov89/go-developer-course-gophkeeper/internal/service"
	"github.com/vstebletsov89/go-developer-course-gophkeeper/internal/service/auth"
	"github.com/vstebletsov89/go-developer-course-gophkeeper/internal/storage/postgres"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"

	"google.golang.org/grpc"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/vstebletsov89/go-developer-course-gophkeeper/internal/config"
	pb "github.com/vstebletsov89/go-developer-course-gophkeeper/internal/proto"
)

// GophkeeperServer represents a structure for gophkeeper service.
type GophkeeperServer struct {
	pb.UnimplementedGophkeeperServer
	service service.Service
}

// NewGophkeeperServer returns an instance of GophkeeperServer.
func NewGophkeeperServer(service service.Service) *GophkeeperServer {
	return &GophkeeperServer{service: service}
}

// AddData adds encrypted private data to the storage.
func (g *GophkeeperServer) AddData(ctx context.Context, request *pb.AddDataRequest) (*pb.AddDataResponse, error) {
	userID := auth.ExtractUserIDFromContext(ctx)

	var response pb.AddDataResponse
	data, err := secure.EncryptPrivateData(request.GetData(), userID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	err = g.service.AddData(ctx, data)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	log.Debug().Msg("Server (AddData): done")
	return &response, nil
}

// GetData gets all related data for current user.
func (g *GophkeeperServer) GetData(ctx context.Context, request *pb.GetDataRequest) (*pb.GetDataResponse, error) {
	log.Debug().Msgf("Server (GetData) request: %v", request)

	userID := auth.ExtractUserIDFromContext(ctx)
	var response pb.GetDataResponse

	data, err := g.service.GetDataByUserID(ctx, userID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	for _, v := range data {
		secret, err := secure.DecryptPrivateData(v)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		response.Data = append(response.Data, secret)
	}

	log.Debug().Msg("Server (GetData): done")
	return &response, nil
}

// DeleteData deletes private data from the storage.
func (g *GophkeeperServer) DeleteData(ctx context.Context, request *pb.DeleteDataRequest) (*pb.DeleteDataResponse, error) {
	var response pb.DeleteDataResponse

	err := g.service.DeleteDataByDataID(ctx, request.GetDataId())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	log.Debug().Msg("Server (DeleteData): done")
	return &response, nil
}

// RunServer starts server application for gophkeeper service.
//
//nolint:funlen
func RunServer(cfg *config.Config) error {
	// init global logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// default log level is info
	zerolog.SetGlobalLevel(config.ParseLogLevel(cfg.LogLevel))

	// debug config
	log.Debug().Msgf("%+v\n\n", cfg)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// error group to control server instances
	g, _ := errgroup.WithContext(ctx)

	db, err := postgres.ConnectDB(context.Background(), cfg.DatabaseDsn)
	if err != nil {
		log.Debug().Msgf("connectDB error: %s", err)
		return err
	}

	log.Info().Msg("Database connection: OK")
	defer db.Close()

	conn, err := postgres.ConnectDBForMigration(cfg.DatabaseDsn)
	if err != nil {
		return err
	}
	log.Info().Msg("Migration connection: OK")

	defer func(conn *sql.DB) {
		err := conn.Close()
		if err != nil {
			return
		}
	}(conn)

	// run migrations
	if cfg.EnableMigration {
		err := postgres.RunMigrations(conn)
		if err != nil {
			return err
		}
	}

	// create storage
	storage := postgres.NewDBStorage(db)

	// create new service
	svc := service.NewService(storage)

	jwtManager := auth.NewJWTManager(cfg.JwtSecretKey)
	jwt := NewJwtInterceptor(jwtManager)
	authServer := NewAuthServer(*svc, jwtManager)
	gophkeeperServer := NewGophkeeperServer(*svc)

	var grpcSrv *grpc.Server

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	// start GRPC server with/without TLS
	g.Go(func() error {
		listen, err := net.Listen("tcp", cfg.ServerAddress)
		if err != nil {
			log.Error().Msgf("GRPC server net.Listen: %v", err.Error())
			return err
		}

		if cfg.EnableTLS {
			// Server using TLS credentials
			log.Info().Msg("GRPC server configuration with TLS credentials")
			transportCredentials, err := credentials.NewServerTLSFromFile("cert.pem", "key.pem")
			if err != nil {
				log.Error().Msgf("GRPC server credentials.NewServerTLSFromFile: %v", err.Error())
				return err
			}

			grpcSrv = grpc.NewServer(
				grpc.Creds(transportCredentials),
				grpc.UnaryInterceptor(jwt.UnaryInterceptor))
		} else {
			// server without TLS credentials
			log.Info().Msg("GRPC server configuration without TLS credentials")
			grpcSrv = grpc.NewServer(
				grpc.UnaryInterceptor(jwt.UnaryInterceptor))
		}

		pb.RegisterAuthServer(grpcSrv, authServer)
		pb.RegisterGophkeeperServer(grpcSrv, gophkeeperServer)

		log.Info().Msgf("GRPC server started on %v", cfg.ServerAddress)
		// start grc server
		if err := grpcSrv.Serve(listen); err != nil {
			log.Error().Msgf("Serve error: %v", err.Error())
			return err
		}
		return nil
	})

	<-sigint

	grpcSrv.GracefulStop()

	// stop server context and release resources
	cancel()

	// release resources
	storage.ReleaseStorage()
	log.Info().Msg("Server Shutdown gracefully")

	err = g.Wait()
	if err != nil {
		log.Error().Msgf("error group: %v", err)
		return err
	}

	return nil
}
