package v4

type PeerGetter interface {
	Get(group string, key string) ([]byte, error)
}
