package main

import "C"

import (
	"encoding/json"
	"fmt"
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
	resMap["ok"] = true

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
	resMap["ok"] = true
	resMap["data"] = "get-resource-return-data: " + string(data)

	res, err := json.Marshal(resMap)
	if err != nil {
		return cgoError(err)
	}

	return C.CString(string(res))
}

// util apis
func cgoError(err error) *C.char {
	return C.CString("Error:: " + err.Error())
}

// main
func main() {}
