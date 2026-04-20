package appgrpc

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	authRPC "github.com/Gusuv/sso/internal/grpc/auth"

	"net"

	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(log *slog.Logger, port int, authService authRPC.AuthService) *App {
	gRPCServer := grpc.NewServer()

	authRPC.Register(gRPCServer, authService)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

const (
	runOp  = "appgrpc.Run"
	stopOp = "appgrpc.Stop"
)

func (a *App) Run() error {

	log := a.log.With(slog.String("op", runOp), slog.Int("port", a.port))

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))

	if err != nil {
		return fmt.Errorf("%s: %w", runOp, err)
	}

	log.Info("gRPC Run", slog.String("addr", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", runOp, err)
	}

	return nil
}

func (a *App) stop() {
	a.gRPCServer.GracefulStop()
	a.log.With(slog.String("op", stopOp)).Info("gRPC Stopped", slog.Int("port", a.port))
}

func (a *App) ReadStop() {

	stopped := make(chan os.Signal, 1)
	signal.Notify(stopped, syscall.SIGTERM, syscall.SIGINT)

	stopName := <-stopped

	a.log.Info("Stopping", slog.String("signal", stopName.String()))
	a.stop()
	a.log.Info("Stopped")

}
