// Copyright (c) 2023 by Alibaba.
// Licensed under the Apache License, Version 2.0, see LICENSE for details.
// SPDX-License-Identifier: Apache-2.0

use super::{Repository, ResourceDesc};
use anyhow::{Result};
use serde::Deserialize;
use std::path::{PathBuf};
use log::info;
use std::process::{Command, Stdio};
use std::io::{self, Read};

pub const DEFAULT_REPO_DIR_PATH: &str = "/opt/confidential-containers/kbs/repository";

#[derive(Debug, Deserialize, Clone)]
pub struct LocalFsRepoDesc {
    pub dir_path: Option<String>,
}

impl Default for LocalFsRepoDesc {
    fn default() -> Self {
        Self {
            dir_path: Some(DEFAULT_REPO_DIR_PATH.to_string()),
        }
    }
}

pub struct LocalFs {
    pub repo_dir_path: String,
    pub cfsi: crate::cfs::Cfs,
}

#[async_trait::async_trait]
impl Repository for LocalFs {
    async fn read_secret_resource(&self, resource_desc: ResourceDesc, extra_request: &str) -> Result<Vec<u8>> {
        let get_res = self.cfsi.get_resource(resource_desc.repository_name.clone(),
                                             resource_desc.resource_type.clone(),
                                             resource_desc.resource_tag.clone(),
                                             extra_request)
            .await?;
        info!("confilesystem - cfsi.get_resource() -> get_res = {:?}", get_res);
        Ok(get_res)

        /*
        let _resource_path = PathBuf::from(&self.repo_dir_path);

        let ref_resource_path = format!(
            "{}/{}/{}",
            resource_desc.repository_name, resource_desc.resource_type, resource_desc.resource_tag
        );
        info!("read resource {}", ref_resource_path);
        let mut output = Command::new("cfs-resource")
            .arg("get")
            .arg("-d")
            .arg(&self.repo_dir_path)
            .arg("-r").arg(resource_desc.repository_name)
            .arg("-k").arg(resource_desc.resource_type)
            .arg("-t").arg(resource_desc.resource_tag)
            .stdout(Stdio::piped())
            .spawn()?;
        let mut stdout = output.stdout.take().expect("Failed to take stdout");
        let mut buffer: Vec<u8> = Vec::new();
        stdout.read_to_end(&mut buffer)?;

        let status = output.wait()?;
        if !status.success() {
            return Err(io::Error::new(io::ErrorKind::Other,
                                      format!("fail to read {}", ref_resource_path),).into());
        }
        Ok(buffer)
        */
        /*
        resource_path.push(ref_resource_path);

        let resource_byte = tokio::fs::read(&resource_path)
            .await
            .context("read resource from local fs")?;
        Ok(resource_byte)
         */
    }
}

impl LocalFs {
    pub fn new(repo_desc: &LocalFsRepoDesc) -> Result<Self> {
        let repo_dir = repo_desc
            .dir_path
            .clone()
            .unwrap_or(DEFAULT_REPO_DIR_PATH.to_string());
        Ok(Self {
            repo_dir_path: repo_dir.clone(),
            cfsi: crate::cfs::Cfs::new("local fs".to_string(), repo_dir.clone())?,
        })
    }
}

#[cfg(test)]
mod tests {
    use crate::resource::{
        local_fs::{LocalFs, LocalFsRepoDesc},
        Repository, ResourceDesc,
    };

    const TEST_DATA: &[u8] = b"testdata";

    #[tokio::test]
    async fn write_and_read_resource() {
        let tmp_dir = tempfile::tempdir().expect("create temp dir failed");
        let repo_desc = LocalFsRepoDesc {
            dir_path: Some(tmp_dir.path().to_string_lossy().to_string()),
        };

        let mut local_fs = LocalFs::new(&repo_desc).expect("create local fs failed");
        let resource_desc = ResourceDesc {
            repository_name: "default".into(),
            resource_type: "test".into(),
            resource_tag: "test".into(),
        };

        let data = local_fs
            .read_secret_resource(resource_desc, "extra-request-test")
            .await
            .expect("read secret resource failed");

        assert_eq!(&data[..], TEST_DATA);
    }
}
