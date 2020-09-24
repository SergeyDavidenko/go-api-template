package storage



var (
	// StorageDB ...
	StorageDB Storage
)


// Storage interface
type Storage interface {
	Init() error
	ShowVersion() string
	Close() error
}
