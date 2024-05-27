package main

import (
	"fmt"
)

const (
	ResourceStorageTypeFile = "file"
	ResourceStorageTypeDb   = "db"
)

type Resource interface {
	SetResource(addr, typ, tag string, data []byte) error
	GetResource(addr, typ, tag string) ([]byte, error)
}

type Config struct {
	StorageType string `json:"storageType,omitempty"`
}

func NewResource(config Config) (Resource, error) {
	switch config.StorageType {
	case ResourceStorageTypeFile:
		return NewFile(config)
	case ResourceStorageTypeDb:
		return nil, fmt.Errorf("ResourceStorageTypeDb: Not implement")
	}

	return nil, fmt.Errorf("Config: Not support")
}
