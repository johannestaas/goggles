package kvstore

import "log"


type KVStore struct {
    db string
    memmap map[string]string
}


func (store *KVStore) Get(key string) string {
    log.Printf("db %s getting key %s\n", store.db, key)
    result := store.memmap[key]
    return result
}

func (store *KVStore) Set(key string, val string) {
    log.Printf("db %s setting key %s to %s\n", store.db, key, val)
    store.memmap[key] = val
}

func New(db string) *KVStore {
    store := new(KVStore)
    store.db = db
    store.memmap = make(map[string]string)
    log.Printf("instanciated %s KVStore\n", db)
    return store
}
