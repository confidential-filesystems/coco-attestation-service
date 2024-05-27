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
    pub fn setResource(path: GoString, data: GoString) -> *mut c_char;
    pub fn getResource(path: GoString) -> *mut c_char;
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
        resource_path: String,
        resource_data: String,
    ) -> Result<bool> {
        let path_go = GoString {
            p: resource_path.as_ptr() as *const c_char,
            n: resource_path.len() as isize,
        };

        let data_go = GoString {
            p: resource_data.as_ptr() as *const c_char,
            n: resource_data.len() as isize,
        };

        log::debug!("confilesystem - set_resource(): resource_path: {:?}, resource_data: {:?}",
            resource_path, resource_data);
        // Call the function exported by cgo and process
        let res_buf: *mut c_char =
            unsafe { setResource(path_go, data_go) };
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

        Ok(result_boolean)
    }

    pub async fn get_resource(
        &self,
        resource_path: String
    ) -> Result<(bool, String)> {
        let path_go = GoString {
            p: resource_path.as_ptr() as *const c_char,
            n: resource_path.len() as isize,
        };

        log::debug!("confilesystem - get_resource(): resource_path: {:?}", resource_path);
        // Call the function exported by cgo and process
        let res_buf: *mut c_char =
            unsafe { getResource(path_go) };
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
        let result_data = res_kv["data"]
            .to_string();
            //.ok_or_else(|| anyhow!("CFS output must contain \"data\" String value"))?;

        Ok((result_boolean, result_data))
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
