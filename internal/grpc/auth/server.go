package auth

import (
	"context"
	"errors"
	"main/internal/service"
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
}

type Auth interface {
	UserLogin(ctx context.Context, email, password string, appID int64) (accessToken, refreshToken string, userId, expiresAt int64, err error)
	UserRegister(ctx context.Context, username, email, password string) (success bool, err error)
	AdminCheck(ctx context.Context, accessToken string) (bool, error)
	LogOut(ctx context.Context, refreshToken string) (string, error)
}

func Register(grpcServ *grpc.Server, auth Auth) {
	sso1.RegisterAuthServiceServer(grpcServ, &serverAPI{
		auth: auth,
	})

}

func (s *serverAPI) Login(ctx context.Context, req *sso1.LoginRequest) (*sso1.LoginResponse, error) {

	email := strings.TrimSpace(req.GetEmail())
	password := req.GetPassword()

	if err := validation.LoginValidation(email, password, req.GetAppId()); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	accessToken, refreshToken, userId, expiresAt, err := s.auth.UserLogin(ctx, email, password, req.GetAppId())
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserNotFound):
			return nil, status.Error(codes.NotFound, "User not found")
		case errors.Is(err, service.ErrInvalidPassword):
			return nil, status.Error(codes.Unauthenticated, "Invalid password")
		default:
			return nil, status.Error(codes.Internal, "Something went wrong")
		}
	}
	return &sso1.LoginResponse{AccessToken: accessToken, RefreshToken: refreshToken, UserId: userId, ExpiresAt: expiresAt}, nil
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
		switch {
		case errors.Is(err, service.ErrUserAlreadyExist):
			return nil, status.Error(codes.AlreadyExists, "User already exists")

		case errors.Is(err, service.ErrUserCreating):
			return nil, status.Error(codes.Internal, "Failed to create user")

		case errors.Is(err, service.ErrPasswordHashing):
			return nil, status.Error(codes.Internal, "Something went wrong with password")

		default:
			return nil, status.Error(codes.Internal, "Something went wrong")

		}

	}

	return &sso1.RegisterResponse{Success: success}, nil
}

func (s *serverAPI) Logout(ctx context.Context, req *sso1.LogoutRequest) (*sso1.LogoutResponse, error) {
	return &sso1.LogoutResponse{Message: "Successful user logout"}, nil
}

func (s *serverAPI) AdminCheck(ctx context.Context, req *sso1.AdminCheckRequest) (*sso1.AdminCheckResponse, error) {
	return &sso1.AdminCheckResponse{IsAdmin: false}, nil
}
