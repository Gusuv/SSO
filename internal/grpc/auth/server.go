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
	// TODO: Add Service methods
	// In development
}

type serverAPI struct {
	sso1.UnimplementedAuthServiceServer
}

func Register(grpc *grpc.Server) {
	sso1.RegisterAuthServiceServer(grpc, &serverAPI{})
}

func requiredData(dataFields map[string]string) error {
	for name, value := range dataFields {
		if strings.TrimSpace(value) == "" {
			return status.Errorf(codes.InvalidArgument, "%s is required", name)
		}
	}
	return nil
}

func (s *serverAPI) Login(ctx context.Context, req *sso1.LoginRequest) (*sso1.LoginResponse, error) {
	email := strings.TrimSpace(req.GetEmail())
	password := strings.TrimSpace(req.GetPassword())

	if err := requiredData(map[string]string{
		"Email":    email,
		"Password": password}); err != nil {
		return nil, err
	}

	return &sso1.LoginResponse{}, nil
}

func (s *serverAPI) Register(ctx context.Context, req *sso1.RegisterRequest) (*sso1.RegisterResponse, error) {
	email := strings.TrimSpace(req.GetEmail())
	password := strings.TrimSpace(req.GetPassword())
	username := strings.TrimSpace(req.GetUsername())

	if err := requiredData(map[string]string{
		"Email":    email,
		"Password": password,
		"Username": username,
	}); err != nil {
		return nil, err
	}
	return &sso1.RegisterResponse{Success: true}, nil
}

func (s *serverAPI) Logout(ctx context.Context, req *sso1.LogoutRequest) (*sso1.LogoutResponse, error) {
	panic("In development")
}

func (s *serverAPI) AdminCheck(ctx context.Context, req *sso1.AdminCheckRequest) (*sso1.AdminCheckResponse, error) {
	panic("In development")
}
