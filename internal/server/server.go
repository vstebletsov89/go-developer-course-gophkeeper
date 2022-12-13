// Package server contains grpc servers and interceptors.
package server

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vstebletsov89/go-developer-course-gophkeeper/internal/secure"
	"github.com/vstebletsov89/go-developer-course-gophkeeper/internal/service"
	"github.com/vstebletsov89/go-developer-course-gophkeeper/internal/service/auth"
	"github.com/vstebletsov89/go-developer-course-gophkeeper/internal/storage/postgres"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/vstebletsov89/go-developer-course-gophkeeper/internal/config"
	pb "github.com/vstebletsov89/go-developer-course-gophkeeper/internal/proto"
)

type GophkeeperServer struct {
	pb.UnimplementedGophkeeperServer
	service service.Service
}

func NewGophkeeperServer(service service.Service) *GophkeeperServer {
	return &GophkeeperServer{service: service}
}

func (g *GophkeeperServer) AddData(ctx context.Context, request *pb.AddDataRequest) (*pb.AddDataResponse, error) {
	var response pb.AddDataResponse
	data, err := secure.EncryptPrivateData(request.GetData())
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

func (g *GophkeeperServer) GetData(ctx context.Context, request *pb.GetDataRequest) (*pb.GetDataResponse, error) {
	var response pb.GetDataResponse

	data, err := g.service.GetDataByUserID(ctx, request.GetUserId())
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

func (g *GophkeeperServer) DeleteData(ctx context.Context, request *pb.DeleteDataRequest) (*pb.DeleteDataResponse, error) {
	var response pb.DeleteDataResponse

	err := g.service.DeleteDataByDataID(ctx, request.GetDataId())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	log.Debug().Msg("Server (DeleteData): done")
	return &response, nil
}

func parseLogLevel(level string) zerolog.Level {
	switch level {
	case "trace":
		return zerolog.TraceLevel
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	case "panic":
		return zerolog.PanicLevel
	case "disabled":
		return zerolog.Disabled
	}
	return zerolog.InfoLevel
}

func RunServer(cfg *config.Config) error {
	// init global logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// default log level is info
	zerolog.SetGlobalLevel(parseLogLevel(cfg.LogLevel))

	// debug config
	log.Debug().Msgf("%+v\n\n", cfg)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// error group to control server instances
	g, ctx := errgroup.WithContext(ctx)

	db, err := connectDB(context.Background(), cfg.DatabaseDsn)
	if err != nil {
		log.Debug().Msgf("connectDB error: %s", err)
		return err
	}

	log.Info().Msg("Database connection: OK")
	defer db.Close()

	// create storage
	storage := postgres.NewDBStorage(db)

	// create new service
	svc := service.NewService(storage)
	var grpcSrv *grpc.Server

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	// TODO: add TLS support and credentials
	// start GRPC server with TLS
	g.Go(func() error {
		listen, err := net.Listen("tcp", cfg.ServerAddress)
		if err != nil {
			log.Error().Msgf("GRPC server net.Listen: %v", err.Error())
			return err
		}

		jwtManager := auth.NewJWTManager(cfg.LogLevel)

		grpcSrv = grpc.NewServer(grpc.UnaryInterceptor(UnaryInterceptor))
		pb.RegisterAuthServer(grpcSrv, NewAuthServer(*svc, jwtManager))
		pb.RegisterGophkeeperServer(grpcSrv, NewGophkeeperServer(*svc))

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

func connectDB(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	log.Debug().Msg("Connect to DB...")
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		log.Error().Msgf("Failed to connect to database. Error: %v", err.Error())
		return nil, err
	}
	return pool, nil
}

func UnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// TODO: think about it (extract current JWT token and check it is valid)

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		values := md.Get("service.AccessToken")
		if len(values) > 0 {
			//userID = values[0]
			//log.Printf("UnaryInterceptor userID from context: '%s'", userID)
		}
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		//md.Append("service.AccessToke"n, string("userID"))
	}
	newCtx := metadata.NewIncomingContext(ctx, md)

	return handler(newCtx, req)
}
