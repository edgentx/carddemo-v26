package shared

// DomainEvent interface
type DomainEvent interface {
	Type() string
}
