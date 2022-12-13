// Package server contains grpc servers and interceptors.
package server

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vstebletsov89/go-developer-course-gophkeeper/internal/secure"
	"github.com/vstebletsov89/go-developer-course-gophkeeper/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	_, cancel := context.WithCancel(context.Background()) //TODO: return ctx
	defer cancel()

	db, err := connectDB(context.Background(), cfg.DatabaseDsn)
	if err != nil {
		log.Debug().Msgf("connectDB error: %s", err)
		return err
	}

	log.Info().Msg("Database connection: OK")
	defer db.Close()

	// create new service
	var grpcSrv *grpc.Server

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	// start GRPC server with TLS
	// TODO: implement it

	<-sigint

	// graceful shutdown
	grpcSrv.GracefulStop()

	// stop server context and release resources
	cancel()

	// release resources
	//err := storage.ReleaseStorage() // TODO: uncomment it
	//if err != nil {
	//	log.Error().Msgf("Release storage error %v", err)
	//}
	log.Info().Msg("Server Shutdown gracefully")

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
