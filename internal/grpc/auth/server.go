package auth

import (
	"context"
	"log/slog"
	"main/internal/validation"
	"strings"

	sso1 "github.com/Gusuv/sso-protos/generated/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverAPI struct {
	sso1.UnimplementedAuthServiceServer
	auth Auth
	log  *slog.Logger
}

type Auth interface {
	UserLogin(ctx context.Context, email, password string, appID int64) (accessToken, refreshToken string, userId int64, err error)
	UserRegister(ctx context.Context, username, email, password string) (success bool, err error)
	AdminCheck(ctx context.Context, accessToken string) (bool, error)
}

func Register(grpcServ *grpc.Server, auth Auth, log *slog.Logger) {
	sso1.RegisterAuthServiceServer(grpcServ, &serverAPI{
		auth: auth,
		log:  log,
	})

}

func (s *serverAPI) Login(ctx context.Context, req *sso1.LoginRequest) (*sso1.LoginResponse, error) {

	email := strings.TrimSpace(req.GetEmail())
	password := req.GetPassword()

	if err := validation.LoginValidation(email, password, req.GetAppId()); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	accessToken, refreshToken, userId, err := s.auth.UserLogin(ctx, email, password, req.GetAppId())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to login")
	}
	return &sso1.LoginResponse{AccessToken: accessToken, RefreshToken: refreshToken, UserId: userId}, nil
}

func (s *serverAPI) Register(ctx context.Context, req *sso1.RegisterRequest) (*sso1.RegisterResponse, error) {
	username := strings.TrimSpace(req.GetUsername())
	email := strings.TrimSpace(req.GetEmail())
	password := req.GetPassword()

	if err := validation.RegisterValidation(username, email, password); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	success, err := s.auth.UserRegister(ctx, username, email, password)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to register user")
	}

	return &sso1.RegisterResponse{Success: success}, nil
}

func (s *serverAPI) Logout(ctx context.Context, req *sso1.LogoutRequest) (*sso1.LogoutResponse, error) {
	return &sso1.LogoutResponse{Message: "Successful user logout"}, nil
}

func (s *serverAPI) AdminCheck(ctx context.Context, req *sso1.AdminCheckRequest) (*sso1.AdminCheckResponse, error) {
	return &sso1.AdminCheckResponse{IsAdmin: true}, nil
}
