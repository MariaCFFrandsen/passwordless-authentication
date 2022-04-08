package db_access

import (
	".authenticator/utils"
	"github.com/dgraph-io/badger"
)

type Access struct {
	Db *badger.DB
}

func Connect(dbPath string) (*Access, *badger.DB) { //receiver?
	var opts badger.Options
	opts = badger.DefaultOptions(dbPath)
	opts.Truncate = true
	db, err := badger.Open(opts)
	utils.Handle(err)
	return &Access{db}, db
}

func (a *Access) Close() {
	a.Db.Close()
}
