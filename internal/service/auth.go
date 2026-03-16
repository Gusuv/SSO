package service

import (
	"context"
	"errors"
	"log/slog"
	"main/internal/repository"
	"main/internal/security/hash"
	"time"
)

type AuthService struct {
	log      *slog.Logger
	rep      AuthRepo
	tokenTTL time.Duration
}

type AuthRepo interface {
	TxCreateUser(ctx context.Context, username, email, passwordHash string) (int64, error)
	// In development
}

func New(log *slog.Logger, rep AuthRepo, tokenTTL time.Duration) *AuthService {
	return &AuthService{log: log, rep: rep, tokenTTL: tokenTTL}
}

func (a *AuthService) UserLogin(ctx context.Context, email, password string, appID int64) (accessToken, refreshToken string, userId int64, err error) {
	return "In development", "In development", 0, nil
}

func (a *AuthService) UserRegister(ctx context.Context, username, email, password string) (success bool, err error) {

	passwordHash, err := hash.MakeHash(password)
	if err != nil {
		a.log.Error("password hashing error", slog.Any("error", err))
		return false, ErrPasswordHashing
	}

	id, err := a.rep.TxCreateUser(ctx, username, email, passwordHash)
	if err != nil {
		if errors.Is(err, repository.UserExist) {
			return false, ErrUserAlreadyExist
		}

		if errors.Is(err, repository.SetRoleError) {
			a.log.Error("register failed: set role error", slog.Any("error", err))
			return false, ErrUserCreating
		}
		a.log.Error("register failed: something went wrong", slog.Any("error", err))
		return false, err
	}

	a.log.Info("user Successful registered",
		slog.Int64("id", id),
		slog.String("email", email))

	return true, nil
}

func (a *AuthService) AdminCheck(ctx context.Context, accessToken string) (bool, error) {
	// In development
	return false, nil
}
