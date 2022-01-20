package mysqldb

import "bitbucket.org/HeilaSystems/persistentstorage"

func NewMySQLDb() persistentstorage.PersistentStorage {
	return MySQLDb{}
}

type MySQLDb struct {
}

func (db MySQLDb) Query() error {
	return nil
}

func (db MySQLDb) Exec() error {
	return nil
}
