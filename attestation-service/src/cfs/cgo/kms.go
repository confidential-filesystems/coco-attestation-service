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

//export initKMS
func initKMS(storageType string, storeFileRepoDir string) *C.char {
	fmt.Printf("confilesystem-go - initKMS(): storageType = %v, storeFileRepoDir = %v\n",
		storageType, storeFileRepoDir)
	var err error = nil
	config = Config{
		StorageType:      storageType,
		StoreFileRepoDir: storeFileRepoDir,
	}
	// FIXME: unknown reason that the StoreFileRepoDir will be changed after set
	// for example:
	// origin StoreFileRepoDir: /opt/confidential-containers/kbs/repository
	// later resourceInstance.RepoDir will be changed to /opa/confidential-containers/kbs/policy.reg
	resourceInstance, err = NewResource(config)
	fmt.Printf("confilesystem-go - initKMS(): NewResource() -> err = %v\n", err)
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

//export setResource
func setResource(storeFileRepoDir, addr, typ, tag string, data string) *C.char {
	fmt.Printf("confilesystem-go - setResource(): addr = %v, typ = %v, tag = %v, data = %v\n",
		addr, typ, tag, data)

	/*
	   // POST ownership/filesystems/:name
	*/
	switch addr {
	case "ownership":
		{
			switch typ {
			case "filesystems":
				{
					fmt.Printf("confilesystem-go - setResource(): -> mintFilesystem(): filesystem-name = %v\n", tag)
					return mintFilesystem(data)
				}
			}
		}
	}

	err := resourceInstance.SetResource(storeFileRepoDir, addr, typ, tag, ([]byte)(data))
	if err != nil {
		return cgoError(err)
	}

	resMap := make(map[string]interface{})
	resMap[ResMapKeyOk] = true
	resMap[ResMapKeyData] = "secret-resource"

	res, err := json.Marshal(resMap)
	if err != nil {
		return cgoError(err)
	}

	return C.CString(string(res))
}

//export deleteResource
func deleteResource(storeFileRepoDir, addr, typ, tag string, data string) *C.char {
	fmt.Printf("confilesystem-go - deleteResource(): addr = %v, typ = %v, tag = %v, data = %v\n",
		addr, typ, tag, data)

	/*
	   // DELETE ownership/filesystems/:name
	*/
	switch addr {
	case "ownership":
		{
			switch typ {
			case "filesystems":
				{
					fmt.Printf("confilesystem-go - deleteResource(): -> burnFilesystem(): filesystem-name = %v\n", tag)
					return burnFilesystem(data)
				}
			}
		}
	}

	err := resourceInstance.DeleteResource(storeFileRepoDir, addr, typ, tag)
	if err != nil {
		return cgoError(err)
	}

	resMap := make(map[string]interface{})
	resMap[ResMapKeyOk] = true
	resMap[ResMapKeyData] = "delete-resource"

	res, err := json.Marshal(resMap)
	if err != nil {
		return cgoError(err)
	}

	return C.CString(string(res))
}

//export getResource
func getResource(storeFileRepoDir, addr, typ, tag, extraRequest string) *C.char {
	fmt.Printf("confilesystem-go - getResource(): addr = %v, typ = %v, tag = %v\n",
		addr, typ, tag)
	fmt.Printf("confilesystem-go - getResource(): extraRequest = %v\n", extraRequest)

	/*
		// GET ownership/filesystems/:name
		// GET ownership/accounts_metatx/:addr
		// GET ownership/configure/.well-known
	*/
	switch addr {
	case "ownership":
		{
			switch typ {
			case "filesystems":
				{
					return getFilesystem(tag)
				}
			case "accounts_metatx":
				{
					return getAccountMetaTx(tag)
				}
			case "configure":
				{
					switch tag {
					case ".well-known":
						{
							return getWellKnownCfg()
						}
					}
				}
			}
		}
	}

	data, err := resourceInstance.GetResource(storeFileRepoDir, addr, typ, tag, extraRequest)
	if err != nil {
		return cgoError(err)
	}

	resMap := make(map[string]interface{})
	resMap[ResMapKeyOk] = true
	resMap[ResMapKeyData] = data // will be base64 encoded

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
	fmt.Printf("confilesystem-go - verifySeeds(): seeds = %v\n", seeds)
	fmt.Printf("confilesystem-go - verifySeeds(): kl = %v\n", kl)
	fmt.Printf("confilesystem-go - verifySeeds(): err = %v\n", err)
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
