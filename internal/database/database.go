package database

import (
	"internal/kvstore"
	"log"
)

type Database struct {
	DataDir string
	Stores  map[string]*kvstore.KVStore
}

func (db *Database) Persist() {
	// not implemented yet
	return
}

func (db *Database) GetOrCreateStore(name *string) *kvstore.KVStore {
	store, ok := db.Stores[*name]
	if !ok {
		log.Printf("creating db %s\n", *name)
		store = kvstore.New(*name)
		db.Stores[*name] = store
	} else {
		log.Printf("switched to db %s\n", *name)
	}
	return store
}

func (db *Database) DropStore(store *kvstore.KVStore) {
	delete(db.Stores, store.Name)
}

func New(dataDir *string) *Database {
	db := new(Database)
	db.DataDir = *dataDir
	db.Stores = make(map[string]*kvstore.KVStore)
	log.Printf("instanciated new database at %s\n", *dataDir)
	return db
}
