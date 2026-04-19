package service

import (
	"context"
	"errors"
	"log/slog"
	"main/internal/models"
	"main/internal/repository"
	security "main/internal/security/jwt"
	"time"
)

type AuthService struct {
	log  *slog.Logger
	rep  AuthRepo
	hash AuthHash
	jwt  AuthToken
}

type AuthRepo interface {
	AddRefreshToken(ctx context.Context, jwt *security.Tokens, refersHash string) error
	GetUserWithRole(ctx context.Context, email string) (*models.Users, []string, error)
	TxCreateUser(ctx context.Context, username, email, passwordHash string) (int64, error)
	CheckRefreshToken(ctx context.Context, refreshToken string) (*models.RefreshTokens, error)
	GetUserRolesById(ctx context.Context, userId int64) ([]string, error)
	RevokeRefreshToken(ctx context.Context, tokenHash string) error
}

const (
	loginOp    = "Auth.login"
	registerOp = "Auth.register"
	logoutOp   = "Auth.logout"
	refreshOp  = "Auth.refresh"
	getMeOp    = "Auth.getMe"
)

type AuthToken interface {
	GenerateTokens(userId int64, role []string) (*security.Tokens, error)
	GenerateAccessToken(userID int64, role []string) (*security.AccessToken, error)
}

type AuthHash interface {
	HashCompare(hash, password string) bool
	HashToken(refreshToken string) string
	MakeHash(password string) (string, error)
}

func New(log *slog.Logger, rep AuthRepo, jwt AuthToken, hash AuthHash) *AuthService {
	return &AuthService{log: log, rep: rep, jwt: jwt, hash: hash}
}

func (a *AuthService) UserLogin(ctx context.Context, email, password string) (LoginResult, error) {
	loginLog := a.log.With("op", loginOp)
	start := time.Now()

	user, role, err := a.rep.GetUserWithRole(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			loginLog.Warn("user not found", slog.Any("err", err))
			return LoginResult{}, ErrInvalidCredentials
		}
		loginLog.Error("failed to get user", slog.Any("err", err))
		return LoginResult{}, err
	}

	if !a.hash.HashCompare(user.PasswordHash, password) {
		loginLog.Warn("invalid password")
		return LoginResult{}, ErrInvalidCredentials
	}

	jwt, err := a.jwt.GenerateTokens(user.Id, role)
	if err != nil {
		loginLog.Error("error with jwt", slog.Any("err", err))
		return LoginResult{}, err
	}
	refreshHash := a.hash.HashToken(jwt.RefreshToken)

	if err := a.rep.AddRefreshToken(ctx, jwt, refreshHash); err != nil {

		loginLog.Error("cant add refresh token to db", slog.Any("err", err))
		return LoginResult{}, err

	}

	result := LoginResult{
		AccessToken:  jwt.AccessToken,
		RefreshToken: jwt.RefreshToken,
		UserId:       user.Id,
		ExpiresAt:    jwt.AccessExpiresAt.Unix(),
	}

	loginLog.Info("login success",
		slog.Int64("user_id", user.Id),
		slog.String("username", user.Username),
		slog.Duration("latency", time.Since(start)),
		slog.Any("role", role),
	)

	return result, nil
}

func (a *AuthService) UserRegister(ctx context.Context, username, email, password string) (userId int64, err error) {
	regLog := a.log.With(slog.String("op", registerOp))
	start := time.Now()

	passwordHash, err := a.hash.MakeHash(password)
	if err != nil {
		regLog.Error("password hashing error", slog.Any("err", err))
		return 0, ErrPasswordHashing
	}

	id, err := a.rep.TxCreateUser(ctx, username, email, passwordHash)
	if err != nil {
		if errors.Is(err, repository.ErrUserExist) {
			return 0, ErrUserAlreadyExist
		}

		if errors.Is(err, repository.ErrSetRoleError) {
			regLog.Error("register failed: set role error", slog.Any("err", err))
			return 0, ErrUserCreating
		}
		regLog.Error("register failed: something went wrong", slog.Any("err", err))
		return 0, err
	}

	regLog.Info("register success",
		slog.Int64("user_id", id),
		slog.String("email", email),
		slog.Duration("latency", time.Since(start)),
	)
	return id, nil
}

func (a *AuthService) Logout(ctx context.Context, refreshToken string) error {
	/*
		outLog := a.log.With(slog.String("op", logoutOp))
		start := time.Now()
	*/

	return nil
}

func (a *AuthService) GetUser(ctx context.Context) (result *GetUserResult, err error) {
	return nil, nil
}

func (a *AuthService) RefreshToken(ctx context.Context, refreshToken string) (result *RefreshResult, err error) {
	refLog := a.log.With(slog.String("op", refreshOp))
	start := time.Now()

	refreshTokenHash := a.hash.HashToken(refreshToken)

	ref, err := a.rep.CheckRefreshToken(ctx, refreshTokenHash)
	if err != nil {
		if errors.Is(err, repository.ErrRefreshTokenNotFound) {
			return nil, ErrInvalidRefreshToken
		}
		return nil, err
	}

	if ref.Revoked {
		refLog.Warn("refresh token is revoked", slog.Int64("user id", ref.UserId))
		return nil, ErrInvalidRefreshToken
	}
	if ref.ExpiresAt.Before(time.Now()) {
		refLog.Warn("refresh token is expired", slog.Int64("user id", ref.UserId))
		err := a.rep.RevokeRefreshToken(ctx, refreshTokenHash)
		if err != nil {
			return nil, err
		}
		return nil, ErrRefreshTokenExpired
	}

	userRoles, err := a.rep.GetUserRolesById(ctx, ref.UserId)
	if err != nil {
		return nil, err
	}

	newAccessToken, err := a.jwt.GenerateAccessToken(ref.UserId, userRoles)
	if err != nil {
		return nil, err
	}

	refLog.Info("refresh token success",
		slog.Duration("latency", time.Since(start)),
		slog.Int64("user_id", ref.UserId),
	)

	return &RefreshResult{
		AccessToken: newAccessToken.AccessToken,
		Refresh:     refreshToken,
		ExpiresAt:   newAccessToken.AccessExpiresAt,
		UserId:      ref.UserId,
	}, nil
}

type GetUserResult struct {
	UserId   int64
	Username string
	Email    string
	Roles    []string
}
type RefreshResult struct {
	AccessToken string
	Refresh     string
	ExpiresAt   time.Time
	UserId      int64
}
type LoginResult struct {
	AccessToken  string
	RefreshToken string
	UserId       int64
	ExpiresAt    int64
}
