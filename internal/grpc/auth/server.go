package auth

import (
	"context"

	sso1 "github.com/Gusuv/sso-protos/generated/go/sso"
	"google.golang.org/grpc"
)

type serverAPI struct {
	sso1.UnimplementedAuthServer
}

func Register(grpc *grpc.Server) {
	sso1.RegisterAuthServer(grpc, &serverAPI{})
}

func (s *serverAPI) Login(ctx context.Context, req *sso1.LoginReq) (*sso1.LoginResp, error) {
	panic("In development")
}

func (s *serverAPI) Register(ctx context.Context, req *sso1.RegisterReq) (*sso1.RegisterResp, error) {
	panic("In development")
}

func (s *serverAPI) Logout(ctx context.Context, req *sso1.LogoutReq) (*sso1.LogoutResp, error) {
	panic("In development")
}

func (s *serverAPI) AdminCheck(ctx context.Context, req *sso1.AdminCheckReq) (*sso1.AdminCheckResp, error) {
	panic("In development")
}
