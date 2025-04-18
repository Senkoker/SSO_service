package kit

import (
	"GRPC_Service_sso/config"
	"context"
	"database/sql"
	"fmt"
	sso_v1_ssov1 "github.com/Senkoker/sso_proto/proto/proto_go/protobufcontract/protobufcontract"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"strconv"
	"testing"
)

type Kit_st struct {
	Ctx         context.Context
	Auth_client sso_v1_ssov1.AuthClient
	Db          *sql.DB
	Cfg         config.ClientParser
}

func Kit_new(t *testing.T) *Kit_st {
	t.Parallel()
	cfg := config.ClientConfigparser()
	fmt.Println(cfg)
	addres := net.JoinHostPort(cfg.Addres, strconv.Itoa(cfg.Port))
	cc, err := grpc.NewClient(addres, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), cfg.Idletimeout)
	client := sso_v1_ssov1.NewAuthClient(cc)
	db, err := sql.Open("pgx", cfg.Database_url)
	if err != nil {
		log.Fatalln(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalln(err)
	}
	return &Kit_st{Ctx: ctx, Auth_client: client, Db: db, Cfg: cfg}
}
