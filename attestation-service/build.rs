use std::process::{exit, Command};

fn real_main() -> Result<(), String> {
    let out_dir = std::env::var("OUT_DIR").unwrap();
    println!("cargo:rerun-if-changed={out_dir}");
    println!("cargo:rustc-link-search=native={out_dir}");
    println!("cargo:rustc-link-lib=static=cgo");
    let cgo_dir = "./src/cgo".to_string();
    let cgo = Command::new("go")
        .args([
            "build",
            "-o",
            &format!("{out_dir}/libcgo.a"),
            "-buildmode=c-archive",
            "opa.go",
            "intoto.go",
        ])
        .current_dir(cgo_dir)
        .output()
        .expect("failed to launch opa compile process");
    if !cgo.status.success() {
        return Err(std::str::from_utf8(&cgo.stderr.to_vec())
            .unwrap()
            .to_string());
    }

    tonic_build::compile_protos("../protos/attestation.proto").map_err(|e| format!("{e}"))?;

    tonic_build::compile_protos("../protos/reference.proto").map_err(|e| format!("{e}"))?;

    {
        let out_dir = std::env::var("OUT_DIR").unwrap();
        println!("cargo:rerun-if-changed={out_dir}");
        println!("cargo:rustc-link-search=native={out_dir}");
        println!("cargo:rustc-link-lib=dylib=cfs"); // static, dylib (the default), or framework
        let cgo_dir = "./src/cfs/cgo".to_string();
        let cgo = Command::new("go")
            .args([
                "build",
                "-o",
                &format!("{out_dir}/libcfs.so"),
                "-buildmode=c-shared", // "-buildmode=c-archive",
                "util.go",
                "ownership.go",
                "resource.go",
                "file.go",
                "cfs.go",
            ])
            .current_dir(cgo_dir)
            .output()
            .expect("failed to launch cfs compile process");
        if !cgo.status.success() {
            return Err(std::str::from_utf8(&cgo.stderr.to_vec())
                .unwrap()
                .to_string());
        }
    }

    Ok(())
}

fn main() -> shadow_rs::SdResult<()> {
    if let Err(e) = real_main() {
        eprintln!("ERROR: {e}");
        exit(1);
    }

    shadow_rs::new()
}
