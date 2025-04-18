package app

import (
	"GRPC_Service_sso/config"
	"GRPC_Service_sso/internal/Auth"
	"GRPC_Service_sso/internal/grpc/grpc_functions"
	storage2 "GRPC_Service_sso/internal/storage"
	"google.golang.org/grpc"
	log2 "log"
	"log/slog"
)

func App(log *slog.Logger, cfg config.Cfg, server *grpc.Server) {
	storage, err := storage2.NewSt(cfg.Database_url)
	if err != nil {
		log2.Fatal("Can't connect to database")
	}
	Auth := Auth.NewAuth(log, cfg, &storage, &storage, &storage, &storage)
	grpc_functions.Server_regist(server, Auth)
}
