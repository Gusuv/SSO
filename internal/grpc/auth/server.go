package auth

import (
	"context"
	"strings"

	sso1 "github.com/Gusuv/sso-protos/generated/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
}

type serverAPI struct {
	sso1.UnimplementedAuthServiceServer
}

func Register(grpc *grpc.Server) {
	sso1.RegisterAuthServiceServer(grpc, &serverAPI{})
}

func (s *serverAPI) Login(ctx context.Context, req *sso1.LoginRequest) (*sso1.LoginResponse, error) {
	email := strings.TrimSpace(req.GetEmail())
	password := strings.TrimSpace(req.GetPassword())

	if email == "" {
		return nil, status.Error(codes.InvalidArgument, "Email is required")
	}
	if password == "" {
		return nil, status.Error(codes.InvalidArgument, "Password is required")
	}
	return &sso1.LoginResponse{}, nil
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
