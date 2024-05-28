// Copyright (c) 2023 by Alibaba.
// Licensed under the Apache License, Version 2.0, see LICENSE for details.
// SPDX-License-Identifier: Apache-2.0

use anyhow::{anyhow, Result};
//use async_trait::async_trait;
use serde_json::Value;
use std::ffi::CStr;
use std::os::raw::c_char;

// Link import cgo function
#[link(name = "cfs")]
extern "C" {
    pub fn setResource(addr: GoString, typ: GoString, tag: GoString, data: GoString) -> *mut c_char;
    pub fn getResource(addr: GoString, typ: GoString, tag: GoString) -> *mut c_char;
    pub fn verifySeeds(seeds: GoString) -> *mut c_char;
}

/// String structure passed into cgo
#[derive(Debug)]
#[repr(C)]
pub struct GoString {
    pub p: *const c_char,
    pub n: isize,
}

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
    pub async fn set_resource(
        &self,
        repository_name: String,
        resource_type: String,
        resource_tag: String,
        resource_data: &[u8],
    ) -> Result<()> {
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

        log::debug!("confilesystem - set_resource(): addr_go: {:?}, typ_go: {:?}, tag_go: {:?}",
            addr_go, typ_go, tag_go);
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

        log::debug!("confilesystem - get_resource(): addr_go: {:?}, typ_go: {:?}, tag_go: {:?}",
            addr_go, typ_go, tag_go);
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
        let seeds_go = GoString {
            p: seeds.as_ptr() as *const c_char,
            n: seeds.len() as isize,
        };

        log::debug!("confilesystem - verify_seeds(): seeds_go: {:?}", seeds_go);
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
}
