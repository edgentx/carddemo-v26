package storage

// StorageClient is an interface for file storage operations (S3/MinIO).
type StorageClient interface {
	GetFile(key string) ([]byte, string, error)
	PutFile(key string, data []byte) (string, error)
}
