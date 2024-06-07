package main

import "C"

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/confidential-filesystems/filesystem-toolchain/resource"
)

var (
	config = Config{
		StorageType: ResourceStorageTypeFile,
	}

	resourceInstance Resource = nil
)

func init() {
	var err error = nil
	resourceInstance, err = NewResource(config)
	fmt.Printf("confilesystem-go - init(): NewResource() -> err = %v\n", err)
	if err != nil {
		panic(err)
	}
}

//export setResource
func setResource(addr, typ, tag string, data string) *C.char {
	//
	fmt.Printf("confilesystem-go - setResource(): addr = %v, typ = %v, tag = %v, data = %v\n",
		addr, typ, tag, data)

	err := resourceInstance.SetResource(addr, typ, tag, ([]byte)(data))
	if err != nil {
		return cgoError(err)
	}

	resMap := make(map[string]interface{})
	resMap[ResMapKeyOk] = true

	res, err := json.Marshal(resMap)
	if err != nil {
		return cgoError(err)
	}

	return C.CString(string(res))
}

//export getResource
func getResource(addr, typ, tag string) *C.char {
	//
	fmt.Printf("confilesystem-go - getResource(): addr = %v, typ = %v, tag = %v\n",
		addr, typ, tag)

	data, err := resourceInstance.GetResource(addr, typ, tag)
	if err != nil {
		return cgoError(err)
	}

	resMap := make(map[string]interface{})
	resMap[ResMapKeyOk] = true
	resMap[ResMapKeyData] = string(data)

	res, err := json.Marshal(resMap)
	if err != nil {
		return cgoError(err)
	}

	return C.CString(string(res))
}

//export verifySeeds
func verifySeeds(seeds string) *C.char {
	if seeds == "" {
		return cgoError(errors.New("seeds is empty"))
	}
	kl, err := resource.NewKeyLoad(seeds)
	if err != nil {
		return cgoError(err)
	}
	if !kl.Valid() {
		return cgoError(errors.New("seeds is invalid"))
	}

	resMap := make(map[string]interface{})
	resMap[ResMapKeyOk] = true

	res, err := json.Marshal(resMap)
	if err != nil {
		return cgoError(err)
	}

	return C.CString(string(res))
}

// main
func main() {}
