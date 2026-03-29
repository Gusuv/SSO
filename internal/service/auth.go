package service

import (
	"context"
	"errors"
	"log/slog"
	"main/internal/models"
	"main/internal/repository"
	"main/internal/security/hash"
	"time"
)

type AuthService struct {
	log      *slog.Logger
	rep      AuthRepo
	jwt      AuthToken
	tokenTTL time.Duration
}

type AuthRepo interface {
	TxCreateUser(ctx context.Context, username, email, passwordHash string) (int64, error)
	GetUserByEmail(ctx context.Context, email string) (*models.Users, *models.Roles, error)
}
type AuthToken interface {
	GenerateToken(userId, appId int64, role string) (string, string, error)
}

func New(log *slog.Logger, rep AuthRepo, tokenTTL time.Duration, jwt AuthToken) *AuthService {
	return &AuthService{log: log, rep: rep, tokenTTL: tokenTTL, jwt: jwt}
}

func (a *AuthService) UserLogin(ctx context.Context, email, password string, appID int64) (accessT string, refreshT string, userId, expiresAt int64, err error) {
	user, role, err := a.rep.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return "", "", 0, 0, ErrUserNotFound
		}
		a.log.Error("failed to get user", slog.Any("error", err))
		return "", "", 0, 0, err
	}

	if !hash.HashCompare(user.PasswordHash, password) {
		a.log.Warn("invalid password", slog.Int64("id", user.Id))
		return "", "", 0, 0, ErrInvalidPassword
	}
	a.log.Info("user successfully logged in",
		slog.Int64("id", user.Id),
		slog.String("username", user.Username))

	accessToken, refreshToken, err := a.jwt.GenerateToken(user.Id, appID, role.Role)
	if err != nil {
		a.log.Error("error with jwt", slog.Any("err", err))
		return "", "", 0, 0, err
	}

	return accessToken, refreshToken, user.Id, time.Now().Add(a.tokenTTL).Unix(), nil
}

func (a *AuthService) UserRegister(ctx context.Context, username, email, password string) (success bool, err error) {

	passwordHash, err := hash.MakeHash(password)
	if err != nil {
		a.log.Error("password hashing error", slog.Any("error", err))
		return false, ErrPasswordHashing
	}

	id, err := a.rep.TxCreateUser(ctx, username, email, passwordHash)
	if err != nil {
		if errors.Is(err, repository.ErrUserExist) {
			return false, ErrUserAlreadyExist
		}

		if errors.Is(err, repository.ErrSetRoleError) {
			a.log.Error("register failed: set role error", slog.Any("error", err))
			return false, ErrUserCreating
		}
		a.log.Error("register failed: something went wrong", slog.Any("error", err))
		return false, err
	}

	a.log.Info("user successful registered",
		slog.Int64("id", id),
		slog.String("email", email))

	return true, nil
}

func (a *AuthService) AdminCheck(ctx context.Context, accessToken string) (bool, error) {
	// In development
	return false, nil
}
