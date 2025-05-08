package modules

import (
	"context"

	"github.com/ashwinath/simple/framework"
	"github.com/ashwinath/simple/server"
)

type serverModule struct {
	fw framework.FW
}

// For pprof
func NewServerModule(fw framework.FW) framework.Module {
	return &serverModule{
		fw: fw,
	}
}

func (m *serverModule) Name() string {
	return "Server"
}

func (m *serverModule) Start(ctx context.Context) {
	server.NewServer(m.fw).Serve()
}
