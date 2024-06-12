// Copyright (c) 2023 by Alibaba.
// Licensed under the Apache License, Version 2.0, see LICENSE for details.
// SPDX-License-Identifier: Apache-2.0

use anyhow::{anyhow, Result};
//use async_trait::async_trait;
use serde_json::Value;
use std::ffi::CStr;
use std::os::raw::c_char;
use kbs_types::TeePubKey;
use serde::{Serialize, Deserialize};
use crate::rvps::store::StoreType;

// Link import cgo function
#[link(name = "cfs")]
extern "C" {
    // <-> kms,ca
    pub fn initKMS(storage_type: GoString) -> *mut c_char;
    pub fn setResource(addr: GoString, typ: GoString, tag: GoString, data: GoString) -> *mut c_char;
    pub fn getResource(addr: GoString, typ: GoString, tag: GoString) -> *mut c_char;
    pub fn verifySeeds(seeds: GoString) -> *mut c_char;

    // <-> ownership
    pub fn initOwnership(cfg_file: GoString, ctx_timeout_sec: i64) -> *mut c_char;
    pub fn mintFilesystem(req: GoString) -> *mut c_char;
    pub fn getFilesystem(name: GoString) -> *mut c_char;
    pub fn burnFilesystem(req: GoString) -> *mut c_char;
    pub fn getAccountMetaTx(addr: GoString) -> *mut c_char;
    pub fn getWellKnownCfg() -> *mut c_char;
}

/// String structure passed into cgo
#[derive(Debug)]
#[repr(C)]
pub struct GoString {
    pub p: *const c_char,
    pub n: isize,
}

// rust to go struct
#[derive(Serialize, Deserialize, Debug)]
pub struct MetaTxForwardRequest {
    #[serde(rename = "from")]
    pub from: String,
    #[serde(rename = "to")]
    pub to: String,
    #[serde(rename = "value")]
    pub value: String,
    #[serde(rename = "gas")]
    pub gas: String,
    #[serde(rename = "nonce")]
    pub nonce: String,
    #[serde(rename = "deadline")]
    pub deadline: u64,
    #[serde(rename = "data")]
    pub data: String,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct MintFilesystemReq {
    #[serde(rename = "metaTxRequest")]
    pub meta_tx_request: MetaTxForwardRequest,
    #[serde(rename = "metaTxSignature")]
    pub meta_tx_signature: String,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct MintFilesystemResp {
    #[serde(rename = "token_id")]
    pub token_id: String,
}

pub type BurnFilesystemReq = MintFilesystemReq;

#[derive(Serialize, Deserialize, Debug)]
pub struct GetFilesystemResp {
    #[serde(rename = "filesystem_name")]
    pub filesystem_name: String,
    #[serde(rename = "owner_address")]
    pub owner_address: String,
    #[serde(rename = "token_id")]
    pub token_id: String,
    #[serde(rename = "token_uri")]
    pub token_uri: String,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct GetMetaTxParamsResp {
    #[serde(rename = "chain_id")]
    pub chain_id: u64,
    #[serde(rename = "chain_number")]
    pub chain_number: u64,
    #[serde(rename = "eip712_name")]
    pub eip712_name: String,
    #[serde(rename = "eip712_version")]
    pub eip712_version: String,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct Contracts {
    #[serde(rename = "forwarder")]
    pub forwarder: String,
    #[serde(rename = "filesystem")]
    pub filesystem: String,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct GetConfigureResp {
    #[serde(rename = "configure")]
    pub configure: GetMetaTxParamsResp,
    #[serde(rename = "contracts")]
    pub contracts: Contracts,
    #[serde(rename = "nonce")]
    pub nonce: u64,
}

// cfs
#[derive(Debug, Clone)]
pub struct Cfs {
    info: String,
}

impl Cfs {
    pub fn new(info: String) -> Result<Self> {

        Ok(Self { info })
    }
}

// #[async_trait]
impl Cfs {
    // init cfs
    pub fn init_cfs(
        kms_store_type: String,
        ownership_cfg_file: String, ownership_ctx_timeout_sec: i64
    ) -> Result<()> {
        // init kms
        log::debug!("confilesystem - init_cfs() - initKMS(): kms_store_type: {:?}", kms_store_type);
        let kms_store_type_go = GoString {
            p: kms_store_type.as_ptr() as *const c_char,
            n: kms_store_type.len() as isize,
        };

        // Call the function exported by cgo and process
        let res_buf: *mut c_char =
            unsafe { initKMS(kms_store_type_go) };
        let res_str: &CStr = unsafe { CStr::from_ptr(res_buf) };
        let res = res_str.to_str()?.to_string();
        log::info!("confilesystem - init_cfs() - initKMS(): res = {:?}", res);
        if res.starts_with("Error::") {
            return Err(anyhow!(res));
        }

        let res_kv: Value = serde_json::from_str(&res)?;
        let result_boolean = res_kv["ok"]
            .as_bool()
            .ok_or_else(|| anyhow!("CFS output must contain \"ok\" boolean value"))?;
        if !result_boolean {
            return Err(anyhow!("CFS output result_boolean is false"));
        }

        // init ownershio
        log::debug!("confilesystem - init_cfs() - initOwnership(): ownership_cfg_file: {:?}, ownership_ctx_timeout_sec: {:?}",
            ownership_cfg_file, ownership_ctx_timeout_sec);
        let ownership_cfg_file_go = GoString {
            p: ownership_cfg_file.as_ptr() as *const c_char,
            n: ownership_cfg_file.len() as isize,
        };

        // Call the function exported by cgo and process
        let res_buf: *mut c_char =
            unsafe { initOwnership(ownership_cfg_file_go, ownership_ctx_timeout_sec) };
        let res_str: &CStr = unsafe { CStr::from_ptr(res_buf) };
        let res = res_str.to_str()?.to_string();
        log::info!("confilesystem - init_cfs() - initOwnership(): res = {:?}", res);
        if res.starts_with("Error::") {
            return Err(anyhow!(res));
        }

        let res_kv: Value = serde_json::from_str(&res)?;
        let result_boolean = res_kv["ok"]
            .as_bool()
            .ok_or_else(|| anyhow!("CFS output must contain \"ok\" boolean value"))?;
        if !result_boolean {
            return Err(anyhow!("CFS output result_boolean is false"));
        }

        Ok(())
    }

    // <-> kms,ca
    pub async fn set_resource(
        &self,
        repository_name: String,
        resource_type: String,
        resource_tag: String,
        resource_data: &[u8],
    ) -> Result<()> {
        log::debug!("confilesystem - set_resource(): repository_name: {:?}, resource_type: {:?}, resource_tag: {:?}",
            repository_name, resource_type, resource_tag);

        let addr_go = GoString {
            p: repository_name.as_ptr() as *const c_char,
            n: repository_name.len() as isize,
        };

        let typ_go = GoString {
            p: resource_type.as_ptr() as *const c_char,
            n: resource_type.len() as isize,
        };

        let tag_go = GoString {
            p: resource_tag.as_ptr() as *const c_char,
            n: resource_tag.len() as isize,
        };

        let data_go = GoString {
            p: resource_data.as_ptr() as *const c_char,
            n: resource_data.len() as isize,
        };

        // Call the function exported by cgo and process
        let res_buf: *mut c_char =
            unsafe { setResource(addr_go, typ_go, tag_go, data_go) };
        let res_str: &CStr = unsafe { CStr::from_ptr(res_buf) };
        let res = res_str.to_str()?.to_string();
        log::info!("confilesystem - set_resource(): res = {:?}", res);
        if res.starts_with("Error::") {
            return Err(anyhow!(res));
        }

        let res_kv: Value = serde_json::from_str(&res)?;
        let result_boolean = res_kv["ok"]
            .as_bool()
            .ok_or_else(|| anyhow!("CFS output must contain \"ok\" boolean value"))?;
        if !result_boolean {
            return Err(anyhow!("CFS output result_boolean is false"));
        }
        Ok(())
    }

    pub async fn get_resource(
        &self,
        repository_name: String,
        resource_type: String,
        resource_tag: String,
    ) -> Result<Vec<u8>> {
        log::debug!("confilesystem - get_resource(): repository_name: {:?}, resource_type: {:?}, resource_tag: {:?}",
            repository_name, resource_type, resource_tag);

        let addr_go = GoString {
            p: repository_name.as_ptr() as *const c_char,
            n: repository_name.len() as isize,
        };

        let typ_go = GoString {
            p: resource_type.as_ptr() as *const c_char,
            n: resource_type.len() as isize,
        };

        let tag_go = GoString {
            p: resource_tag.as_ptr() as *const c_char,
            n: resource_tag.len() as isize,
        };

        // Call the function exported by cgo and process
        let res_buf: *mut c_char =
            unsafe { getResource(addr_go, typ_go, tag_go) };
        let res_str: &CStr = unsafe { CStr::from_ptr(res_buf) };
        let res = res_str.to_str()?.to_string();
        log::info!("confilesystem - get_resource(): res = {:?}", res);
        if res.starts_with("Error::") {
            return Err(anyhow!(res));
        }

        let res_kv: Value = serde_json::from_str(&res)?;
        let result_boolean = res_kv["ok"]
            .as_bool()
            .ok_or_else(|| anyhow!("CFS output must contain \"ok\" boolean value"))?;
        if !result_boolean {
            return Err(anyhow!("CFS output result_boolean is false"));
        }
        let result_data = res_kv["data"]
            .to_string();
            //.ok_or_else(|| anyhow!("CFS output must contain \"data\" String value"))?;

        let result_data_bytes = result_data.into_bytes();
        Ok(result_data_bytes)
    }

    pub fn verify_seeds(
        &self,
        seeds: String,
    ) -> Result<()> {
        log::debug!("confilesystem - verify_seeds(): seeds: {:?}", seeds);

        let seeds_go = GoString {
            p: seeds.as_ptr() as *const c_char,
            n: seeds.len() as isize,
        };

        // Call the function exported by cgo and process
        let res_buf: *mut c_char =
            unsafe { verifySeeds(seeds_go) };
        let res_str: &CStr = unsafe { CStr::from_ptr(res_buf) };
        let res = res_str.to_str()?.to_string();
        log::info!("confilesystem - verify_seeds(): res = {:?}", res);
        if res.starts_with("Error::") {
            return Err(anyhow!(res));
        }

        let res_kv: Value = serde_json::from_str(&res)?;
        let result_boolean = res_kv["ok"]
            .as_bool()
            .ok_or_else(|| anyhow!("CFS output must contain \"ok\" boolean value"))?;
        if !result_boolean {
            return Err(anyhow!("CFS output result_boolean is false"));
        }
        Ok(())
    }

    // <-> ownership
    pub async fn mint_filesystem(
        &self,
        req: &MintFilesystemReq,
    ) -> Result<MintFilesystemResp> {
        log::debug!("confilesystem - mint_filesystem(): req: {:?}", req);

        let req_string = match serde_json::to_string(req) {
            Ok(req_string) => {
                req_string
            },
            Err(e) => {
                anyhow::bail!(e);
            }
        };

        let req_string_go = GoString {
            p: req_string.as_ptr() as *const c_char,
            n: req_string.len() as isize,
        };

        // Call the function exported by cgo and process
        let res_buf: *mut c_char =
            unsafe { mintFilesystem(req_string_go) };
        let res_str: &CStr = unsafe { CStr::from_ptr(res_buf) };
        let res = res_str.to_str()?.to_string();
        log::info!("confilesystem - mint_filesystem(): res = {:?}", res);
        if res.starts_with("Error::") {
            return Err(anyhow!(res));
        }

        let res_kv: Value = serde_json::from_str(&res)?;
        let result_boolean = res_kv["ok"]
            .as_bool()
            .ok_or_else(|| anyhow!("CFS output must contain \"ok\" boolean value"))?;
        if !result_boolean {
            return Err(anyhow!("CFS output result_boolean is false"));
        }
        let result_data = res_kv["data"]
            .to_string();
        //.ok_or_else(|| anyhow!("CFS output must contain \"data\" String value"))?;

        let rsp = serde_json::from_str::<MintFilesystemResp>(&result_data)?;
        Ok(rsp)
    }

    pub async fn get_filesystem(
        &self,
        name: &str,
    ) -> Result<GetFilesystemResp> {
        log::debug!("confilesystem - get_filesystem(): name: {:?}", name);

        let name_go = GoString {
            p: name.as_ptr() as *const c_char,
            n: name.len() as isize,
        };

        // Call the function exported by cgo and process
        let res_buf: *mut c_char =
            unsafe { getFilesystem(name_go) };
        let res_str: &CStr = unsafe { CStr::from_ptr(res_buf) };
        let res = res_str.to_str()?.to_string();
        log::info!("confilesystem - get_filesystem(): res = {:?}", res);
        if res.starts_with("Error::") {
            return Err(anyhow!(res));
        }

        let res_kv: Value = serde_json::from_str(&res)?;
        let result_boolean = res_kv["ok"]
            .as_bool()
            .ok_or_else(|| anyhow!("CFS output must contain \"ok\" boolean value"))?;
        if !result_boolean {
            return Err(anyhow!("CFS output result_boolean is false"));
        }
        let result_data = res_kv["data"]
            .to_string();
        //.ok_or_else(|| anyhow!("CFS output must contain \"data\" String value"))?;

        let rsp = serde_json::from_str::<GetFilesystemResp>(&result_data)?;
        Ok(rsp)
    }

    pub async fn burn_filesystem(
        &self,
        req: &BurnFilesystemReq,
    ) -> Result<()> {
        log::debug!("confilesystem - burn_filesystem(): req: {:?}", req);

        let req_string = match serde_json::to_string(req) {
            Ok(req_string) => {
                req_string
            },
            Err(e) => {
                anyhow::bail!(e);
            }
        };

        let req_string_go = GoString {
            p: req_string.as_ptr() as *const c_char,
            n: req_string.len() as isize,
        };

        // Call the function exported by cgo and process
        let res_buf: *mut c_char =
            unsafe { burnFilesystem(req_string_go) };
        let res_str: &CStr = unsafe { CStr::from_ptr(res_buf) };
        let res = res_str.to_str()?.to_string();
        log::info!("confilesystem - burn_filesystem(): res = {:?}", res);
        if res.starts_with("Error::") {
            return Err(anyhow!(res));
        }

        let res_kv: Value = serde_json::from_str(&res)?;
        let result_boolean = res_kv["ok"]
            .as_bool()
            .ok_or_else(|| anyhow!("CFS output must contain \"ok\" boolean value"))?;
        if !result_boolean {
            return Err(anyhow!("CFS output result_boolean is false"));
        }

        Ok(())
    }

    pub async fn get_account_metatx(
        &self,
        addr: &str,
    ) -> Result<GetMetaTxParamsResp> {
        log::debug!("confilesystem - get_account_metatx(): addr: {:?}", addr);

        let addr_go = GoString {
            p: addr.as_ptr() as *const c_char,
            n: addr.len() as isize,
        };

        // Call the function exported by cgo and process
        let res_buf: *mut c_char =
            unsafe { getAccountMetaTx(addr_go) };
        let res_str: &CStr = unsafe { CStr::from_ptr(res_buf) };
        let res = res_str.to_str()?.to_string();
        log::info!("confilesystem - get_account_metatx(): res = {:?}", res);
        if res.starts_with("Error::") {
            return Err(anyhow!(res));
        }

        let res_kv: Value = serde_json::from_str(&res)?;
        let result_boolean = res_kv["ok"]
            .as_bool()
            .ok_or_else(|| anyhow!("CFS output must contain \"ok\" boolean value"))?;
        if !result_boolean {
            return Err(anyhow!("CFS output result_boolean is false"));
        }
        let result_data = res_kv["data"]
            .to_string();
        //.ok_or_else(|| anyhow!("CFS output must contain \"data\" String value"))?;

        let rsp = serde_json::from_str::<GetMetaTxParamsResp>(&result_data)?;
        Ok(rsp)
    }

    pub async fn get_wellknown(
        &self,
    ) -> Result<GetConfigureResp> {
        log::debug!("confilesystem - get_wellknown():");

        // Call the function exported by cgo and process
        let res_buf: *mut c_char =
            unsafe { getWellKnownCfg() };
        let res_str: &CStr = unsafe { CStr::from_ptr(res_buf) };
        let res = res_str.to_str()?.to_string();
        log::info!("confilesystem - get_wellknown(): res = {:?}", res);
        if res.starts_with("Error::") {
            return Err(anyhow!(res));
        }

        let res_kv: Value = serde_json::from_str(&res)?;
        let result_boolean = res_kv["ok"]
            .as_bool()
            .ok_or_else(|| anyhow!("CFS output must contain \"ok\" boolean value"))?;
        if !result_boolean {
            return Err(anyhow!("CFS output result_boolean is false"));
        }
        let result_data = res_kv["data"]
            .to_string();
        //.ok_or_else(|| anyhow!("CFS output must contain \"data\" String value"))?;

        let rsp = serde_json::from_str::<GetConfigureResp>(&result_data)?;
        Ok(rsp)
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[tokio::test]
    async fn test_set_resource() {

    }

    #[tokio::test]
    async fn test_get_policy() {

    }

    //...
}
