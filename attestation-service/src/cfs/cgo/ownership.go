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

var (
	gOwnershipInitErrStub error = nil

	gCfgFile       string = ""
	gCtxTimeoutSec int64  = 0
)

//export initOwnership
func initOwnership(cfgFile string, ctxTimeoutSec int64) *C.char {
	gOwnershipInitErrStub = nil
	gCfgFile = cfgFile
	gCtxTimeoutSec = ctxTimeoutSec

	err := utils.InitOwnerServFunc(cfgFile, ctxTimeoutSec)
	if err != nil {
		gOwnershipInitErrStub = err
		return cgoError(err)
	}

	resMap := make(map[string]interface{})
	resMap[ResMapKeyOk] = true

	res, err := json.Marshal(resMap)
	if err != nil {
		gOwnershipInitErrStub = err
		return cgoError(err)
	}

	gOwnershipInitErrStub = nil
	return C.CString(string(res))
}

// mint-filesystem
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
	MetaTxRequest   MetaTxForwardRequest `json:"meta_tx_request"`
	MetaTxSignature string               `json:"meta_tx_signature"`
	Meta            string               `json:"meta"`
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
	fmt.Printf("confilesystem-go - mintFilesystem(): req = %v\n", req)

	goReq := MintFilesystemReq{}
	err := json.Unmarshal(([]byte)(req), &goReq)
	if err != nil {
		return cgoError(err)
	}

	mintFilesystemReq, err := toMintFilesystemReq(&goReq)
	if err != nil {
		return cgoError(err)
	}

	if gOwnershipInitErrStub != nil {
		fmt.Printf("mintFilesystem() In: but gOwnershipInitErrStub = %v -> reInit", gOwnershipInitErrStub)
		err = utils.InitOwnerServFunc(gCfgFile, gCtxTimeoutSec)
		if err != nil {
			gOwnershipInitErrStub = err
			return cgoError(err)
		}
		gOwnershipInitErrStub = nil
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

// get-filesystem
type GetFilesystemResp struct {
	FilesystemName string `json:"filesystem_name"`
	OwnerAddress   string `json:"owner_address"`
	TokenId        string `json:"token_id"`
	TokenUri       string `json:"token_uri"`
	Meta           string `json:"meta"`
}

func toGetFilesystemRsp(getFilesystemResp *response.GetFilesystemResp) (*GetFilesystemResp, error) {
	rsp := &GetFilesystemResp{
		FilesystemName: getFilesystemResp.FilesystemName,
		OwnerAddress:   getFilesystemResp.OwnerAddress,
		TokenId:        getFilesystemResp.TokenId,
		TokenUri:       getFilesystemResp.TokenUri,
		Meta:           getFilesystemResp.Meta,
	}
	return rsp, nil
}

//export getFilesystem
func getFilesystem(name string) *C.char {
	fmt.Printf("confilesystem-go - getFilesystem(): name = %v\n", name)

	if gOwnershipInitErrStub != nil {
		fmt.Printf("getFilesystem() In: but gOwnershipInitErrStub = %v -> reInit", gOwnershipInitErrStub)
		err := utils.InitOwnerServFunc(gCfgFile, gCtxTimeoutSec)
		if err != nil {
			gOwnershipInitErrStub = err
			return cgoError(err)
		}
		gOwnershipInitErrStub = nil
	}

	err, getFilesystemResp := utils.GetFs(name)
	if err != nil {
		return cgoError(err)
	}
	rsp, err := toGetFilesystemRsp(getFilesystemResp)
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

// burn-filesystem
type BurnFilesystemReq = MintFilesystemReq

func toBurnFilesystemReq(goReq *BurnFilesystemReq) (*request.BurnFilesystemReq, error) {
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
	req := &request.BurnFilesystemReq{
		MetaTxRequest:   metaTxRequest,
		MetaTxSignature: goReq.MetaTxSignature,
	}

	return req, nil
}

//export burnFilesystem
func burnFilesystem(req string) *C.char {
	fmt.Printf("confilesystem-go - burnFilesystem(): req = %v\n", req)

	goReq := BurnFilesystemReq{}
	err := json.Unmarshal(([]byte)(req), &goReq)
	if err != nil {
		return cgoError(err)
	}

	burnFilesystemReq, err := toBurnFilesystemReq(&goReq)
	if err != nil {
		return cgoError(err)
	}

	if gOwnershipInitErrStub != nil {
		fmt.Printf("burnFilesystem() In: but gOwnershipInitErrStub = %v -> reInit", gOwnershipInitErrStub)
		err = utils.InitOwnerServFunc(gCfgFile, gCtxTimeoutSec)
		if err != nil {
			gOwnershipInitErrStub = err
			return cgoError(err)
		}
		gOwnershipInitErrStub = nil
	}

	err = utils.BurnFs(burnFilesystemReq)
	if err != nil {
		return cgoError(err)
	}

	resMap := make(map[string]interface{})
	resMap[ResMapKeyOk] = true
	resMap[ResMapKeyData] = "burn-filesystem"

	res, err := json.Marshal(resMap)
	if err != nil {
		return cgoError(err)
	}

	return C.CString(string(res))
}

// get-account-metaTx
type GetConfigureResp struct {
	ChainId       uint64 `json:"chain_id"`
	ChainNumber   uint64 `json:"chain_number"`
	EIP712Name    string `json:"eip712_name"`
	EIP712Version string `json:"eip712_version"`
}

type Contracts struct {
	Forwarder  string `json:"forwarder"`
	Filesystem string `json:"filesystem"`
}

type GetMetaTxParamsResp struct {
	GetConfigureResp `json:"configure" mapstructure:"configure"`
	Contracts        `json:"contracts" mapstructure:"contracts"`
	Nonce            uint64 `json:"nonce"`
}

func toGetMetaTxParamsRsp(getMetaTxParamsResp *response.GetMetaTxParamsResp) (*GetMetaTxParamsResp, error) {
	configureResp := GetConfigureResp{
		ChainId:       getMetaTxParamsResp.GetConfigureResp.ChainId,
		ChainNumber:   getMetaTxParamsResp.GetConfigureResp.ChainNumber,
		EIP712Name:    getMetaTxParamsResp.GetConfigureResp.EIP712Name,
		EIP712Version: getMetaTxParamsResp.GetConfigureResp.EIP712Version,
	}

	contracts := Contracts{
		Forwarder:  getMetaTxParamsResp.Contracts.Forwarder,
		Filesystem: getMetaTxParamsResp.Contracts.Filesystem,
	}

	rsp := &GetMetaTxParamsResp{
		GetConfigureResp: configureResp,
		Contracts:        contracts,
		Nonce:            getMetaTxParamsResp.Nonce,
	}
	return rsp, nil
}

//export getAccountMetaTx
func getAccountMetaTx(addr string) *C.char {
	fmt.Printf("confilesystem-go - getAccountMetaTx(): addr = %v\n", addr)

	if gOwnershipInitErrStub != nil {
		fmt.Printf("getAccountMetaTx() In: but gOwnershipInitErrStub = %v -> reInit", gOwnershipInitErrStub)
		err := utils.InitOwnerServFunc(gCfgFile, gCtxTimeoutSec)
		if err != nil {
			gOwnershipInitErrStub = err
			return cgoError(err)
		}
		gOwnershipInitErrStub = nil
	}

	err, getMetaTxParamsResp := utils.GetAccountMetaTx(addr)
	if err != nil {
		return cgoError(err)
	}
	rsp, err := toGetMetaTxParamsRsp(getMetaTxParamsResp)
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

// get-wellknown-cfg
func toGetConfigureRsp(getConfigureResp *response.GetConfigureResp) (*GetConfigureResp, error) {
	rsp := &GetConfigureResp{
		ChainId:       getConfigureResp.ChainId,
		ChainNumber:   getConfigureResp.ChainNumber,
		EIP712Name:    getConfigureResp.EIP712Name,
		EIP712Version: getConfigureResp.EIP712Version,
	}

	return rsp, nil
}

//export getWellKnownCfg
func getWellKnownCfg() *C.char {
	fmt.Printf("confilesystem-go - getWellKnownCfg(): gOwnershipInitErrStub = %v\n", gOwnershipInitErrStub)

	getWellKnownCfgResp := utils.GetWellKnownCfg()
	rsp, err := toGetConfigureRsp(getWellKnownCfgResp)
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
	fmt.Printf("confilesystem-go - getWellKnownCfg(): string(rspBytes) = %v\n", string(rspBytes))

	res, err := json.Marshal(resMap)
	if err != nil {
		return cgoError(err)
	}

	fmt.Printf("confilesystem-go - getWellKnownCfg(): string(res) = %v\n", string(res))
	return C.CString(string(res))
}

//*/
