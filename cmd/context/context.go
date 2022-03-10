package context

import (
	"context"
	"github.com/optman/p2p-tun/host"

	logging "github.com/ipfs/go-log/v2"
)

type loggerKey struct{}
type hostReadyKey struct{}
type clientKey struct{}
type serverKey struct{}
type nodeConfigKey struct{}

type Context struct {
	context.Context
}

func SetLogger(ctx context.Context, v logging.StandardLogger) context.Context {
	return context.WithValue(ctx, loggerKey{}, v)
}

func (ctx *Context) Logger() logging.StandardLogger {
	return ctx.Value(loggerKey{}).(logging.StandardLogger)
}

func SetHostReady(ctx context.Context, v chan struct{}) context.Context {
	return context.WithValue(ctx, hostReadyKey{}, v)
}

func (ctx *Context) HostReady() chan struct{} {
	return ctx.Value(hostReadyKey{}).(chan struct{})
}

func SetClient(ctx context.Context, v *host.Client) context.Context {
	return context.WithValue(ctx, clientKey{}, v)
}

func (ctx *Context) Client() *host.Client {
	return ctx.Value(clientKey{}).(*host.Client)
}

func SetServer(ctx context.Context, v *host.Server) context.Context {
	return context.WithValue(ctx, serverKey{}, v)
}

func (ctx *Context) Server() *host.Server {
	return ctx.Value(serverKey{}).(*host.Server)
}

func SetNodeConfig(ctx context.Context, v *host.NodeConfig) context.Context {
	return context.WithValue(ctx, nodeConfigKey{}, v)
}

func (ctx *Context) NodeConfig() *host.NodeConfig {
	return ctx.Value(nodeConfigKey{}).(*host.NodeConfig)
}
