package server

import (
	"google.golang.org/grpc"
	"log/slog"
	"net"
	"strconv"
)

type Server_info struct {
	Log    *slog.Logger
	Port   int
	Server *grpc.Server
}

func NewServer(log *slog.Logger, port int) Server_info {
	server := grpc.NewServer()
	return Server_info{Log: log, Port: port, Server: server}
}

func (s *Server_info) Start() {
	const op = "Server.Start"
	lis, err := net.Listen("tcp", ":"+strconv.Itoa(s.Port))
	if err != nil {
		s.Log.Error("Cant start server", "error", err, "operation:", op, "port:", strconv.Itoa(s.Port))
	}
	s.Log.Info("server started", op, "port:", strconv.Itoa(s.Port))
	err = s.Server.Serve(lis)
	if err != nil {
		s.Log.Error("Cant start server", "error", err, "operation:", op, "port:", strconv.Itoa(s.Port))
	}
	s.Log.Info("server started", op, "port:", strconv.Itoa(s.Port))
}
func (s *Server_info) Stop() {
	const op = "Server.Stop"
	s.Server.GracefulStop()
	s.Log.Info("Server stopped", op)
}
