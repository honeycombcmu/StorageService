package models

import (
	"errors"
)

type libstore struct {
	storage map[string]interface{}
}

func NewLibStore() (LibStore, error) {
	ls := new(libstore)
	ls.storage = make(map[string]interface{})
	return ls, nil
}

// Get key value pair into database
func (ls *libstore) Get(key string) (interface{}, error) {
	val, ok := ls.storage[key]
	if !ok {
		return nil, errors.New("KeyNotFound")
	}
	return val, nil
}

// Put key value pair into database
func (ls *libstore) Put(key string, val interface{}) error {
	ls.storage[key] = val
	return nil
}
