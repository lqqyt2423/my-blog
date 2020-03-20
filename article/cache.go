package article

import (
	"log"

	"github.com/tidwall/tinylru"
)

// Cache struct
type Cache struct {
	fs  *Fs
	lru *tinylru.LRU
}

// Init func
func (cah *Cache) Init(arg interface{}) {
	cah.fs = new(Fs)
	cah.fs.Init(arg)
	cah.lru = new(tinylru.LRU)
	cah.lru.Resize(10)
}

// Get func
func (cah *Cache) Get(name string) ([]byte, error) {
	key := "Get-" + name

	if res, ok := cah.lru.Get(key); ok {
		log.Printf("Cache Get %v ok\n", name)
		return res.([]byte), nil
	}

	res, err := cah.fs.Get(name)
	if err != nil {
		return nil, err
	}
	cah.lru.Set(key, res)
	return res, nil
}

// GetAll func
func (cah *Cache) GetAll() ([]byte, error) {
	key := "GetAll"
	if res, ok := cah.lru.Get(key); ok {
		log.Printf("Cache GetAll ok\n")
		return res.([]byte), nil
	}

	res, err := cah.fs.GetAll()
	if err != nil {
		return nil, err
	}
	cah.lru.Set(key, res)
	return res, nil
}

// Search func
func (cah *Cache) Search(q string) ([]byte, error) {
	key := "Search-" + q
	if res, ok := cah.lru.Get(key); ok {
		log.Printf("Cache Search %v ok\n", q)
		return res.([]byte), nil
	}

	res, err := cah.fs.Search(q)
	if err != nil {
		return nil, err
	}
	cah.lru.Set(key, res)
	return res, nil
}
