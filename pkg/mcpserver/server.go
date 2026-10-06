// Package mcpserver exposes moneybags data to Model Context Protocol clients
// such as AI assistants, as a set of callable tools.
package mcpserver

import (
	"fmt"
	"net/http"

	"github.com/ashwinath/moneybags/pkg/db"
	"github.com/ashwinath/simple/framework"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const serverName = "moneybags"

// serverVersion is reported to clients during initialisation.
const serverVersion = "0.1.0"

const errorFormatDBNotRegistered = "database %q is not registered or has an unexpected type"

// Server serves the moneybags tools over MCP.
type Server struct {
	fw     framework.FW
	db     db.TransactionDB
	server *mcp.Server
}

// NewServer builds an MCP server with every moneybags tool registered on it.
func NewServer(fw framework.FW) (*Server, error) {
	transactionDB, ok := fw.GetDB(db.TransactionDatabaseName).(db.TransactionDB)
	if !ok {
		return nil, fmt.Errorf(errorFormatDBNotRegistered, db.TransactionDatabaseName)
	}

	s := &Server{
		fw: fw,
		db: transactionDB,
		server: mcp.NewServer(&mcp.Implementation{
			Name:    serverName,
			Version: serverVersion,
		}, nil),
	}

	s.registerTransactionTools()

	return s, nil
}

// HTTPHandler serves the MCP protocol over streamable HTTP.
//
// Sessions are stateless, so each request is served on a throwaway session and
// clients need no affinity across moneybags replicas. It also means the
// handler never accumulates sessions for clients that disconnect without
// cleaning up.
func (s *Server) HTTPHandler() http.Handler {
	return mcp.NewStreamableHTTPHandler(
		func(_ *http.Request) *mcp.Server { return s.server },
		&mcp.StreamableHTTPOptions{
			Stateless: true,
		},
	)
}
