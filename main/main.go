package main

import (
	"GRPC_Service_sso/config"
	"GRPC_Service_sso/internal/app"
	"GRPC_Service_sso/internal/grpc/server"
	"GRPC_Service_sso/internal/logger"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.Cfg_parser()
	log := logger.Newlogger(cfg.Env)
	log.Info("Config", "cfg", cfg)
	Server := server.NewServer(log, cfg.Server.Port)
	app.App(log, cfg, Server.Server)
	go Server.Start()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop
	Server.Stop()
}
