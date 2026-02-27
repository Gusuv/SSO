package auth

import (
	"context"
	"main/internal/validation"
	"strings"

	sso1 "github.com/Gusuv/sso-protos/generated/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	UserLogin(ctx context.Context, email, password string, appID int) (accessToken, refreshToken string, err error)
	UserRegister(ctx context.Context, username, email, password string) (success bool, err error)
	AdminCheck(ctx context.Context, accessToken string) (bool, error)
}

type serverAPI struct {
	sso1.UnimplementedAuthServiceServer
	auth Auth
}

func Register(grpcServ *grpc.Server, auth Auth) {
	sso1.RegisterAuthServiceServer(grpcServ, &serverAPI{
		auth: auth,
	})

}

func (s *serverAPI) Login(ctx context.Context, req *sso1.LoginRequest) (*sso1.LoginResponse, error) {

	email := strings.TrimSpace(req.GetEmail())
	password := strings.TrimSpace(req.GetPassword())
	appid := req.GetAppId()

	if err := validation.LoginValidation(email, password, appid); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &sso1.LoginResponse{}, nil
}

func (s *serverAPI) Register(ctx context.Context, req *sso1.RegisterRequest) (*sso1.RegisterResponse, error) {
	username := strings.TrimSpace(req.GetUsername())
	email := strings.TrimSpace(req.GetEmail())
	password := strings.TrimSpace(req.GetPassword())

	if err := validation.RegisterValidation(username, email, password); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &sso1.RegisterResponse{Success: true}, nil
}

func (s *serverAPI) Logout(ctx context.Context, req *sso1.LogoutRequest) (*sso1.LogoutResponse, error) {
	panic("In development")
}

func (s *serverAPI) AdminCheck(ctx context.Context, req *sso1.AdminCheckRequest) (*sso1.AdminCheckResponse, error) {
	panic("In development")
}
