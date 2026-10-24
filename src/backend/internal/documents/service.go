package documents

// Service owns business rules and transaction boundaries for document requests and staff decisions.
type Service struct {
	store *Store
}

func NewService(store *Store) *Service {
	return &Service{store: store}
}

// TODO: Implement the feature's operations here; keep HTTP details in handler.go.
