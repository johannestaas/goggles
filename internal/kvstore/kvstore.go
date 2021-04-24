package kvstore

import "log"


type KVStore struct {
    db string
    memmap map[string]string
}


func (store *KVStore) Get(key string) string {
    log.Println("db %s getting key %s", store.db, key)
    result := store.memmap[key]
    return result
}

func (store *KVStore) Set(key string, val string) {
    log.Println("db %s setting key %s to %s", store.db, key, val)
    store.memmap[key] = val
}

func New(db string) *KVStore {
    store := new(KVStore)
    store.db = db
    store.memmap = make(map[string]string)
    log.Println("instanciated %s KVStore", db)
    return store
}
