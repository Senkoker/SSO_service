package grpc_functions

import (
	"GRPC_Service_sso/internal/DB_err"
	"context"
	"errors"
	"fmt"
	sso_v1_ssov1 "github.com/Senkoker/sso_proto/proto/proto_go/protobufcontract/protobufcontract"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"regexp"
)

type Server_api struct {
	Server sso_v1_ssov1.UnimplementedAuthServer
	sso_v1_ssov1.UnsafeAuthServer
	Auth Auth
}

func Server_regist(server *grpc.Server, Auth Auth) {
	sso_v1_ssov1.RegisterAuthServer(server, &Server_api{Auth: Auth})
}

// Todo: исправить Protoc
type Auth interface {
	Auth_login(ctx context.Context, email, pass, appid string) (string, error)
	Auth_register(ctx context.Context, email, pass string) (string, error)
	Auth_Accept(ctx context.Context, usercode string) (bool, error)
	Auth_Retry(ctx context.Context, email string) (bool, error)
	Auth_change_pass(ctx context.Context, email, pass string) (bool, error)
	Auth_IsAdmin(ctx context.Context, email string) (bool, error)
}

var (
	Badrequest = errors.New("bad request")
)

func Validation_user(email, password string) error {
	if email == "" || password == "" {
		return Badrequest
	}
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	match, _ := regexp.MatchString(emailRegex, email)
	if match == false {
		return Badrequest
	}
	return nil
}
func (s *Server_api) Login(ctx context.Context, req *sso_v1_ssov1.Loginrequest) (*sso_v1_ssov1.Loginresponse, error) {
	email := req.GetEmail()
	pass := req.GetPassword()
	appid := req.GetAppid()
	err := Validation_user(email, pass)
	if err != nil || appid == "" {
		return nil, status.Error(codes.InvalidArgument, "Email or Password is empty")
	}
	token, err := s.Auth.Auth_login(ctx, email, pass, appid)
	if err != nil {
		fmt.Println(errors.Is(err, DB_err.Does_not_exist))
		if errors.Is(err, DB_err.Does_not_exist) {
			return nil, status.Error(codes.InvalidArgument, "This user is not exist")
		}
		if errors.Is(err, DB_err.Invalid_password) {
			return nil, status.Error(codes.InvalidArgument, "Invalid email or password")
		}
		return nil, status.Error(codes.Internal, "Problem in server")
	}
	return &sso_v1_ssov1.Loginresponse{Token: token}, nil
}

func (s *Server_api) Register(ctx context.Context, req *sso_v1_ssov1.Registrequest) (*sso_v1_ssov1.Registresponse, error) {
	email := req.GetEmail()
	pass := req.GetPassword()
	err := Validation_user(email, pass)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Password or login is empty")
	}
	id, err := s.Auth.Auth_register(ctx, email, pass)
	fmt.Println(errors.Unwrap(err))
	if err != nil {
		if errors.Is(err, DB_err.Dublicate_name) {
			return nil, status.Error(codes.AlreadyExists, "This user already get code in his email")
		}
		if errors.Is(err, DB_err.Already_exists) {
			return nil, status.Error(codes.AlreadyExists, "This user already exists in main db")
		}
		return nil, status.Error(codes.Internal, "Problem in server")
	}
	return &sso_v1_ssov1.Registresponse{Userid: id}, nil

}

func (s *Server_api) ChangePassword(ctx context.Context, req *sso_v1_ssov1.PassChangeRequest) (*sso_v1_ssov1.PassChangeResponse, error) {
	email := req.GetEmail()
	pass := req.GetNewPass()
	if err := Validation_user(email, pass); err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid email or password")
	}
	resp, err := s.Auth.Auth_change_pass(ctx, email, pass)
	if err != nil {
		if errors.Is(err, DB_err.Does_not_exist) {
			return nil, status.Error(codes.InvalidArgument, "This user is not exist")
		}
		if errors.Is(err, DB_err.Dublicate_name) {
			return nil, status.Error(codes.AlreadyExists, "This user already get code in his email")
		}

	}
	return &sso_v1_ssov1.PassChangeResponse{Resp: resp}, nil

}
func (s *Server_api) Accept(ctx context.Context, req *sso_v1_ssov1.Acceptrequest) (*sso_v1_ssov1.Acceptresponse, error) {
	Usercode := req.GetUsercode()
	if Usercode == "" {
		return nil, status.Error(codes.InvalidArgument, "Usercode is empty")
	}
	resp, err := s.Auth.Auth_Accept(ctx, Usercode)
	if err != nil {
		if errors.Is(err, DB_err.Does_not_exist) {
			return nil, status.Error(codes.AlreadyExists, "This user already accept his data")
		}
		return nil, status.Error(codes.Internal, "Problem in server")
	}
	return &sso_v1_ssov1.Acceptresponse{Accresp: resp}, nil

}
func (s *Server_api) Retry(ctx context.Context, req *sso_v1_ssov1.Retryrequest) (*sso_v1_ssov1.Retryresponse, error) {
	user_email := req.GetEmail()
	if user_email == "" {
		return nil, status.Error(codes.InvalidArgument, "Email is empty")
	}
	resp, err := s.Auth.Auth_Retry(ctx, user_email)
	if err != nil {
		if errors.Is(err, DB_err.Invalid_usercode) {
			return nil, status.Error(codes.InvalidArgument, "Invalid usercode")
		}
		return nil, status.Error(codes.Internal, "Problem in server")
	}
	return &sso_v1_ssov1.Retryresponse{Retryresp: resp}, nil

}
func (s *Server_api) IsAdmin(ctx context.Context, req *sso_v1_ssov1.IsAdminrequest) (*sso_v1_ssov1.IsAdminresponse, error) {
	email := req.GetEmail()
	if email == "" {
		return nil, status.Error(codes.InvalidArgument, "Email is empty")
	}
	resp, err := s.Auth.Auth_IsAdmin(ctx, email)
	if err != nil {
		if errors.Is(err, DB_err.Does_not_exist) {
			return nil, status.Error(codes.InvalidArgument, "Not admin")
		}
		return nil, status.Error(codes.Internal, "Problem in server")
	}
	return &sso_v1_ssov1.IsAdminresponse{Adminresp: resp}, nil
}
