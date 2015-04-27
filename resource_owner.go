package oguth

type ResourceOwner interface {
	GetClient(id string) Client
	Close()
}
