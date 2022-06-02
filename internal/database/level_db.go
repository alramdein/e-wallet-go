package database

import (
	"log"

	"github.com/lovoo/goka/storage"
	"github.com/syndtr/goleveldb/leveldb"
)

var (
	Ls storage.Storage
)

func InitLevelDB() {
	DB, err := leveldb.OpenFile("../internal/database/level_db", nil)
	if err != nil {
		log.Fatalf("error instantiating leveldb: %v", err)
	}

	Ls, err = storage.New(DB)
	if err != nil {
		log.Fatalf("error instantiating leveldb: %v", err)
	}
}

func GetLevelDBInstance() storage.Storage {
	return Ls
}
