package internal

import (
	"context"
	"github.com/BlakeLiAFK/ps/pkg"
)

type Context struct {
	Context context.Context
	Auth    pkg.Auth
	PS      *PubSubContext
	server  *httpServer
}

func (s *Context) Run(addr string) {
	srv := newHttpServer(s)
	s.server = srv
	srv.Run(addr)
}
