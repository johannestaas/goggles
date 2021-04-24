package kvstore

import "log"

type KVStore struct {
	name   string
	memmap map[string]string
}

func (store *KVStore) Get(key string) string {
	log.Printf("db %s getting key %s\n", store.name, key)
	result := store.memmap[key]
	return result
}

func (store *KVStore) Set(key string, val string) {
	log.Printf("db %s setting key %s to %s\n", store.name, key, val)
	store.memmap[key] = val
}

func New(name string) *KVStore {
	store := new(KVStore)
	store.name = name
	store.memmap = make(map[string]string)
	log.Printf("instanciated %s KVStore\n", name)
	return store
}
