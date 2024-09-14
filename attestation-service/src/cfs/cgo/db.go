package main

import (
	"context"
	"fmt"
	"time"

	"github.com/confidential-filesystems/filesystem-ownership/utils"
)

const (
	defaultCtxTimeout = time.Second * 15
)

type DB struct {
	Resource
}

func NewDB(_ Config) (Resource, error) {
	db := &DB{}
	return db, nil
}

func (d *DB) SetResource(_, addr, typ, tag string, data []byte) error {
	fmt.Printf("confilesystem-go - DB.SetResource(): addr = %v, typ = %v, tag = %v\n", addr, typ, tag)
	ctx, cancel := context.WithTimeout(context.Background(), defaultCtxTimeout)
	defer cancel()
	return utils.SetResource(ctx, addr, typ, tag, data)
}

func (d *DB) DeleteResource(_, addr, typ, tag string) error {
	fmt.Printf("confilesystem-go - DB.DeleteResource(): addr = %v, typ = %v, tag = %v\n", addr, typ, tag)
	ctx, cancel := context.WithTimeout(context.Background(), defaultCtxTimeout)
	defer cancel()
	return utils.DeleteResource(ctx, addr, typ, tag)
}

func (d *DB) GetResource(_, addr, typ, tag, extraRequest string) ([]byte, error) {
	fmt.Printf("confilesystem-go - DB.GetResource(): addr = %v, typ = %v, tag = %v, extraRequest = %v\n",
		addr, typ, tag, extraRequest)
	ctx, cancel := context.WithTimeout(context.Background(), defaultCtxTimeout)
	defer cancel()
	return utils.ToGetResource(ctx, "", addr, typ, tag, extraRequest)
}
