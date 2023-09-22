package domain

// Source interface
type Source interface {
	GetAllProtos() ([]Proto, error)
	DownloadProto(url string) (string, error)
}
