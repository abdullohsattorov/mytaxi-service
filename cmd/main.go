package main

import (
	"net"

	"google.golang.org/grpc"

	"github.com/abdullohsattorov/mytaxi-service/config"
	pb "github.com/abdullohsattorov/mytaxi-service/genproto"
	"github.com/abdullohsattorov/mytaxi-service/pkg/db"
	"github.com/abdullohsattorov/mytaxi-service/pkg/logger"
	"github.com/abdullohsattorov/mytaxi-service/service"
	"github.com/abdullohsattorov/mytaxi-service/storage"

	"google.golang.org/grpc/reflection"
)

func main() {
	cfg := config.Load()

	log := logger.New(cfg.LogLevel, "myTaxi-service")
	defer func(l logger.Logger) {
		err := logger.Cleanup(l)
		if err != nil {
			log.Fatal("failed cleanup logger", logger.Error(err))
		}
	}(log)

	log.Info("main: sqlxConfig",
		logger.String("host", cfg.PostgresHost),
		logger.Int("port", cfg.PostgresPort),
		logger.String("database", cfg.PostgresDatabase))

	connDB, err := db.ConnectToDB(cfg)
	if err != nil {
		log.Fatal("sqlx connection to postgres error", logger.Error(err))
	}

	pgStorage := storage.NewStoragePg(connDB)

	taxiService := service.NewMyTaxiService(pgStorage, log)

	lis, err := net.Listen("tcp", cfg.RPCPort)
	if err != nil {
		log.Fatal("Error while listening: %v", logger.Error(err))
	}

	s := grpc.NewServer()
	pb.RegisterMyTaxiServiceServer(s, taxiService)
	reflection.Register(s)
	log.Info("main: server running",
		logger.String("port", cfg.RPCPort))

	if err := s.Serve(lis); err != nil {
		log.Fatal("Error while listening: %v", logger.Error(err))
	}
}
