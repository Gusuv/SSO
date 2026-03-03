package service

import (
	"log/slog"
	"time"
)

type AuthService struct {
	log      *slog.Logger
	rep      AuthRepo
	tokenTTL time.Duration
}

type AuthRepo interface {
	createUser(username, email, passwordHash string) error
}

func New(log *slog.Logger, rep AuthRepo, tokenTTL time.Duration) *AuthService {
	return &AuthService{log: log, rep: rep, tokenTTL: tokenTTL}
}
