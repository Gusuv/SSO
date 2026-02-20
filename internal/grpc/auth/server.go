package auth

import (
	"context"

	sso1 "github.com/Gusuv/sso-protos/generated/go/sso"
	"google.golang.org/grpc"
)

type Auth interface {
}

type serverAPI struct {
	sso1.UnimplementedAuthServiceServer
}

func Register(grpc *grpc.Server) {
	sso1.RegisterAuthServiceServer(grpc, &serverAPI{})
}

func (s *serverAPI) Login(ctx context.Context, eq *sso1.LoginRequest) (*sso1.LoginResponse, error) {
	panic("In development")
}

func (s *serverAPI) Register(ctx context.Context, req *sso1.RegisterRequest) (*sso1.RegisterResponse, error) {
	panic("In development")
}

func (s *serverAPI) Logout(ctx context.Context, req *sso1.LogoutRequest) (*sso1.LogoutResponse, error) {
	panic("In development")
}

func (s *serverAPI) AdminCheck(ctx context.Context, req *sso1.AdminCheckRequest) (*sso1.AdminCheckResponse, error) {
	panic("In development")
}
