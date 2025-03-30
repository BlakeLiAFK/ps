package pkg

type Auth interface {
	VerifyToken(token, namespace, topic string) (bool, error) //
}
