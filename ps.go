package ps

import (
	"context"
	"github.com/BlakeLiAFK/ps/internal"
	"github.com/BlakeLiAFK/ps/pkg"
)

func NewAuth(token string) pkg.Auth {
	return &internal.Auth{
		BearerToken: token,
	}
}
func NewServer(auth pkg.Auth) pkg.Server {
	return &internal.Context{
		Context: context.Background(),
		Auth:    auth,
		PS:      internal.NewPubSubContext(),
	}
}
