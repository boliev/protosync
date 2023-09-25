package domain

// Source interface
type Source interface {
	SyncProtos() error
}
