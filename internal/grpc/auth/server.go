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
