package persistentstorage

type PersistentStorage interface {
	Query() error
	Exec() error
}
