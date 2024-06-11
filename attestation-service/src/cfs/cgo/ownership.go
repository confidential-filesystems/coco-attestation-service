package main

import "C"

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/confidential-filesystems/filesystem-ownership/utils"
	"github.com/confidential-filesystems/filesystem-ownership/vos/v1/request"
	"github.com/confidential-filesystems/filesystem-ownership/vos/v1/response"
	"github.com/ethereum/go-ethereum/common"
)

//export initOwnership
func initOwnership(cfgFile string, ctxTimeoutSec int64) *C.char {
	err := utils.InitOwnerServFunc(cfgFile, ctxTimeoutSec)
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

// /*
type MetaTxForwardRequest struct {
	From     string `json:"from"` // common.Address hex
	To       string `json:"to"`   // common.Address hex
	Value    string `json:"value"`
	Gas      string `json:"gas"`
	Nonce    string `json:"nonce"`
	Deadline uint64 `json:"deadline"`
	Data     string `json:"data"` // hex or base64
}

type MintFilesystemReq struct {
	MetaTxRequest   MetaTxForwardRequest `json:"metaTxRequest"`
	MetaTxSignature string               `json:"metaTxSignature"`
}

func toMintFilesystemReq(goReq *MintFilesystemReq) (*request.MintFilesystemReq, error) {
	value, ok := new(big.Int).SetString(goReq.MetaTxRequest.Value, 10)
	if !ok {
		return nil, fmt.Errorf("value %v to big.Int failed", goReq.MetaTxRequest.Value)
	}

	gas, ok := new(big.Int).SetString(goReq.MetaTxRequest.Gas, 10)
	if !ok {
		return nil, fmt.Errorf("gas %v to big.Int failed", goReq.MetaTxRequest.Gas)
	}

	nonce, ok := new(big.Int).SetString(goReq.MetaTxRequest.Nonce, 10)
	if !ok {
		return nil, fmt.Errorf("nonce %v to big.Int failed", goReq.MetaTxRequest.Nonce)
	}

	metaTxRequest := request.MetaTxForwardRequest{
		From:     common.HexToAddress(goReq.MetaTxRequest.From),
		To:       common.HexToAddress(goReq.MetaTxRequest.To),
		Value:    value,
		Gas:      gas,
		Nonce:    nonce,
		Deadline: goReq.MetaTxRequest.Deadline,
		Data:     goReq.MetaTxRequest.Data,
	}
	req := &request.MintFilesystemReq{
		MetaTxRequest:   metaTxRequest,
		MetaTxSignature: goReq.MetaTxSignature,
	}

	return req, nil
}

type MintFilesystemResp struct {
	TokenId string `json:"token_id"`
}

func toMintFilesystemRsp(mintFilesystemResp *response.MintFilesystemResp) (*MintFilesystemResp, error) {
	rsp := &MintFilesystemResp{
		TokenId: mintFilesystemResp.TokenId,
	}
	return rsp, nil
}

//export mintFilesystem
func mintFilesystem(req string) *C.char {
	goReq := MintFilesystemReq{}
	err := json.Unmarshal(([]byte)(req), &goReq)
	if err != nil {
		return cgoError(err)
	}

	mintFilesystemReq, err := toMintFilesystemReq(&goReq)
	if err != nil {
		return cgoError(err)
	}
	err, mintFilesystemResp := utils.MintFs(mintFilesystemReq)
	if err != nil {
		return cgoError(err)
	}
	rsp, err := toMintFilesystemRsp(mintFilesystemResp)
	if err != nil {
		return cgoError(err)
	}

	rspBytes, err := json.Marshal(rsp)
	if err != nil {
		return cgoError(err)
	}

	resMap := make(map[string]interface{})
	resMap[ResMapKeyOk] = true
	resMap[ResMapKeyData] = string(rspBytes)

	res, err := json.Marshal(resMap)
	if err != nil {
		return cgoError(err)
	}

	return C.CString(string(res))
}

//*/
