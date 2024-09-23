package main

import (
	"fmt"
)

const (
	ResourceStorageTypeFile = "file"
	ResourceStorageTypeDb   = "db"
)

type Resource interface {
	SetResource(repoDir, addr, typ, tag string, data []byte) error
	DeleteResource(repoDir, addr, typ, tag string, extraRequest string) error
	GetResource(repoDir, addr, typ, tag, extraRequest string) ([]byte, error)
}

type Config struct {
	StorageType      string `json:"storageType,omitempty"`
	StoreFileRepoDir string `json:"storeFileRepoDir,omitempty"`
}

func NewResource(config Config) (Resource, error) {
	switch config.StorageType {
	case ResourceStorageTypeFile:
		return NewFile(config)
	case ResourceStorageTypeDb:
		return NewDB(config)
	}

	return nil, fmt.Errorf("Config: Not support")
}
