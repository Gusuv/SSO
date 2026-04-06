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
	AddRefreshToken(ctx context.Context, jwt *models.JWT) error
	GetUserWithRole(ctx context.Context, email string) (*models.Users, string, error)
	TxCreateUser(ctx context.Context, username, email, passwordHash string) (int64, error)
}

const (
	loginOp    = "Auth.login"
	registerOp = "Auth.register"
	logoutOp   = "Auth.logout"
	adminOp    = "Auth.admin"
)

type AuthToken interface {
	GenerateToken(userId, appId int64, role string) (*models.JWT, error)
}

func New(log *slog.Logger, rep AuthRepo, tokenTTL time.Duration, jwt AuthToken) *AuthService {
	return &AuthService{log: log, rep: rep, tokenTTL: tokenTTL, jwt: jwt}
}

func (a *AuthService) UserLogin(ctx context.Context, email, password string, appID int64) (accessT string, refreshT string, userId, expiresAt int64, err error) {
	loginLog := a.log.With("op", loginOp, "app_id", appID)
	start := time.Now()

	user, role, err := a.rep.GetUserWithRole(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			loginLog.Warn("user not found", slog.Any("err", err))
			return "", "", 0, 0, ErrUserNotFound
		}
		loginLog.Error("failed to get user", slog.Any("err", err))
		return "", "", 0, 0, err
	}

	if !hash.HashCompare(user.PasswordHash, password) {
		loginLog.Warn("invalid password")
		return "", "", 0, 0, ErrInvalidPassword
	}

	jwt, err := a.jwt.GenerateToken(user.Id, appID, role)
	if err != nil {
		loginLog.Error("error with jwt", slog.Any("err", err))
		return "", "", 0, 0, err
	}
	refreshHash, err := hash.MakeHash(jwt.RefreshToken)

	if err != nil {
		loginLog.Error("refresh token hashing error", slog.Any("err", err))
		return "", "", 0, 0, err
	}
	jwt.RefreshHash = refreshHash

	if err := a.rep.AddRefreshToken(ctx, jwt); err != nil {

		loginLog.Error("cant add refresh token to db", slog.Any("err", err))
		return "", "", 0, 0, err

	}

	loginLog.Info("login success",
		slog.Int64("user_id", user.Id),
		slog.String("username", user.Username),
		slog.Duration("latency", time.Since(start)),
	)

	return jwt.AccessToken, jwt.RefreshToken, user.Id, time.Now().Add(a.tokenTTL).Unix(), nil
}

func (a *AuthService) UserRegister(ctx context.Context, username, email, password string) (success bool, err error) {
	regLog := a.log.With(slog.String("op", registerOp))
	start := time.Now()

	passwordHash, err := hash.MakeHash(password)
	if err != nil {
		regLog.Error("password hashing error", slog.Any("err", err))
		return false, ErrPasswordHashing
	}

	id, err := a.rep.TxCreateUser(ctx, username, email, passwordHash)
	if err != nil {
		if errors.Is(err, repository.ErrUserExist) {
			return false, ErrUserAlreadyExist
		}

		if errors.Is(err, repository.ErrSetRoleError) {
			regLog.Error("register failed: set role error", slog.Any("err", err))
			return false, ErrUserCreating
		}
		regLog.Error("register failed: something went wrong", slog.Any("err", err))
		return false, err
	}

	regLog.Info("register success",
		slog.Int64("user_id", id),
		slog.String("email", email),
		slog.Duration("latency", time.Since(start)),
	)
	return true, nil
}

func (a *AuthService) AdminCheck(ctx context.Context, accessToken string) (bool, error) {
	/*  In development
	admLog := a.log.With(slog.String("op", adminOp))
	start := time.Now()
	*/

	return false, nil
}

func (a *AuthService) LogOut(ctx context.Context, refreshToken string) (string, error) {
	/*   In development
	 logOutLog := a.log.With(slog.String("op", logoutOp))
	start := time.Now()
	*/

	return "", nil
}
