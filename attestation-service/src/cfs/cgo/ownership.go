package main

import "C"

import (
	"github.com/confidential-filesystems/filesystem-ownership/utils"
)

//export initOwneship
func initOwneship(cfgFile string, ctxTimeoutSec int64) *C.char {
	err := utils.InitOwnerServFunc(cfgFile, ctxTimeoutSec)
	return cgoError(err)
}

/*
func mintFilesystem() *C.char {
	func MintFs(req *request.MintFilesystemReq) (error, *response.MintFilesystemResp)
}
*/
