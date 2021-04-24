package kvstore

import (
	"log"
	"time"
)

type Record struct {
	Value string
	Expiry *time.Time
}

type KVStore struct {
	Name   string
	MemMap map[string]*Record
}

func (store *KVStore) Get(key string) string {
	log.Printf("db %s: getting key %s\n", store.Name, key)
	record, ok := store.MemMap[key]
	if !ok {
		return ""
	}
	if record.Expiry == nil {
		return record.Value
	}
	if time.Now().Before(*record.Expiry) {
		return record.Value
	} else {
		return "EXPIRED"
	}
}

func (store *KVStore) Set(key string, val string, duration time.Duration) {
	var expiry_ptr *time.Time
	log.Printf("db %s: setting key %s to %s %s\n", store.Name, key, val, duration)
	if duration == 0 {
		expiry_ptr = nil
	} else {
		expiry := time.Now().Add(duration)
		expiry_ptr = &expiry
	}
	record := Record{
		Value: val,
		Expiry: expiry_ptr,
	}
	store.MemMap[key] = &record
}

func New(name string) *KVStore {
	store := new(KVStore)
	store.Name = name
	store.MemMap = make(map[string]*Record)
	log.Printf("instanciated %s KVStore\n", name)
	return store
}
