package modules

import (
	"context"
	"time"

	"github.com/ashwinath/moneybags/pkg/mcpserver"

	"github.com/ashwinath/simple/framework"
	"github.com/ashwinath/simple/server"
)

// mcpRoute is the path that the MCP protocol is served on.
const mcpRoute = "/mcp"

// serverWriteTimeout is the HTTP write timeout. It is disabled because MCP
// holds streamable HTTP responses open for the duration of a tool call, so the
// 5s default would cut off every streamed response.
const serverWriteTimeout = 0

// serverReadTimeout is the HTTP read timeout. It is raised from the 5s default
// so that a slow client is not disconnected part way through a request.
const serverReadTimeout = 30 * time.Second

type serverModule struct {
	fw        framework.FW
	mcpServer *mcpserver.Server
}

// Serves pprof alongside the MCP endpoint.
func NewServerModule(fw framework.FW) (framework.Module, error) {
	mcpServer, err := mcpserver.NewServer(fw)
	if err != nil {
		return nil, err
	}

	return &serverModule{
		fw:        fw,
		mcpServer: mcpServer,
	}, nil
}

func (m *serverModule) Name() string {
	return "Server"
}

func (m *serverModule) Start(ctx context.Context) {
	server.NewServer(m.fw,
		server.WithWriteTimeout(serverWriteTimeout),
		server.WithReadTimeout(serverReadTimeout),
	).
		RegisterHandler(mcpRoute, m.mcpServer.HTTPHandler()).
		Serve()
}
