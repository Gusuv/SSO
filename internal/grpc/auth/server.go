package auth

import (
	"context"
	"errors"
	"main/internal/models"
	"main/internal/service"
	"main/internal/validation"
	"strings"

	sso1 "github.com/Gusuv/sso-protos/generated/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type serverAPI struct {
	sso1.UnimplementedAuthServiceServer
	auth AuthService
}

type AuthService interface {
	UserLogin(ctx context.Context, email, password string) (loginRes service.LoginResult, err error)
	UserRegister(ctx context.Context, username, email, password string) (userId int64, err error)
	Logout(ctx context.Context, refreshToken string) (string, error)
	RefreshToken(ctx context.Context, refreshToken string) (result *service.RefreshResult, err error)
	GetMe(context.Context) (models.Users, string, error)
}

func Register(grpcServ *grpc.Server, auth AuthService) {
	sso1.RegisterAuthServiceServer(grpcServ, &serverAPI{
		auth: auth,
	})

}

func (s *serverAPI) Login(ctx context.Context, req *sso1.LoginRequest) (*sso1.LoginResponse, error) {

	email := strings.TrimSpace(req.GetEmail())
	password := req.GetPassword()

	if err := validation.LoginValidation(email, password); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	res, err := s.auth.UserLogin(ctx, email, password)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidCredentials):
			return nil, status.Error(codes.Unauthenticated, "Invalid credentials")
		default:
			return nil, status.Error(codes.Internal, "Something went wrong")
		}
	}
	return &sso1.LoginResponse{
		Tokens: &sso1.TokensPair{
			AccessToken:  res.AccessToken,
			RefreshToken: res.RefreshToken,
			ExpiresAt: &timestamppb.Timestamp{
				Seconds: res.ExpiresAt,
			},
		},
		UserId: res.UserId,
	}, nil
}

func (s *serverAPI) Register(ctx context.Context, req *sso1.RegisterRequest) (*sso1.RegisterResponse, error) {
	username := strings.TrimSpace(req.GetUsername())
	email := strings.TrimSpace(req.GetEmail())
	password := req.GetPassword()

	if err := validation.RegisterValidation(username, email, password); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	userId, err := s.auth.UserRegister(ctx, username, email, password)
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

	return &sso1.RegisterResponse{UserId: userId}, nil
}

func (s *serverAPI) Logout(ctx context.Context, req *sso1.LogoutRequest) (*sso1.LogoutResponse, error) {
	refreshToken := strings.TrimSpace(req.GetRefreshToken())
	if refreshToken == "" {
		return nil, status.Error(codes.InvalidArgument, "Refresh token is required")
	}

	return &sso1.LogoutResponse{}, nil
}

func (s *serverAPI) Refresh(ctx context.Context, req *sso1.RefreshRequest) (*sso1.RefreshResponse, error) {
	refreshToken := strings.TrimSpace(req.GetRefreshToken())
	if refreshToken == "" {
		return nil, status.Error(codes.InvalidArgument, "Refresh token is required")
	}
	tokens, err := s.auth.RefreshToken(ctx, refreshToken)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidRefreshToken):
			return nil, status.Error(codes.InvalidArgument, "Invalid refresh token")
		case errors.Is(err, service.ErrRefreshTokenExpired):
			return nil, status.Error(codes.Unauthenticated, "Refresh token expired")
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}

	}

	return &sso1.RefreshResponse{
		Tokens: &sso1.TokensPair{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.Refresh,
			ExpiresAt: &timestamppb.Timestamp{
				Seconds: tokens.ExpiresAt.Unix(),
			},
		}, UserId: tokens.UserId,
	}, nil
}

func (s *serverAPI) GetMe(ctx context.Context, _ *emptypb.Empty) (*sso1.GetMeResponse, error) {
	return &sso1.GetMeResponse{}, nil
}
