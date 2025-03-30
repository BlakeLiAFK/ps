package internal

type Auth struct {
	BearerToken string
}

func (a *Auth) VerifyToken(token, namespace, topic string) (bool, error) {
	return a.BearerToken == token, nil
}
