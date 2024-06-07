package main

import "C"

const (
	ResMapKeyOk   = "ok"
	ResMapKeyData = "data"
)

// util apis
func cgoError(err error) *C.char {
	return C.CString("Error:: " + err.Error())
}
