use crate::modules::state::Files;
use std::collections::HashMap;
use std::fs::{create_dir, read_to_string, File};
use std::io::{self, Write};
use tempdir::TempDir;

use cargo::core::Workspace;
use cargo::util::Config;

pub async fn update_dependencies(
    _req: dependagot_common::UpdateDependenciesRequest,
    files: Files,
) -> Result<impl warp::Reply, warp::Rejection> {
    // Write files out to a temporary directory:
    let sandbox = match setup_sandbox(files).await {
        Err(e) => {
            error!("error creating sandbox: {:?}", e);
            // TODO: custom error
            return Err(warp::reject::not_found());
        }
        Ok(s) => s,
    };
    info!("sandbox: {}", sandbox.path().to_str().unwrap());

    // Parse files into a cargo workspace:
    let config = match Config::default() {
        Err(e) => {
            error!("error initializing cargo config: {:?}", e);
            // TODO: custom error
            return Err(warp::reject::not_found());
        }
        Ok(c) => c,
    };
    let ws = match Workspace::new(&sandbox.path().join("Cargo.toml"), &config) {
        Err(e) => {
            error!("error initializing workspace: {:?}", e);
            // TODO: custom error
            return Err(warp::reject::not_found());
        }
        Ok(s) => s,
    };

    let res = match cargo::ops::update_lockfile(
        &ws,
        &cargo::ops::UpdateOptions {
            config: &config,
            aggressive: false,
            dry_run: false,
            precise: Some("0.6.0"),
            to_update: vec!["prost:0.6.1".to_string()],
        },
    ) {
        Err(e) => {
            error!("error updating dependencies: {:?}", e);
            // TODO: custom error
            return Err(warp::reject::not_found());
        }
        Ok(s) => s,
    };
    info!("res: {:?}", res);

    let new_files = match read_new_files(sandbox) {
        Err(e) => {
            error!("error reading new files: {:?}", e);
            // TODO: custom error
            return Err(warp::reject::not_found());
        }
        Ok(f) => f,
    };

    let res = dependagot_common::UpdateDependenciesResponse { new_files };
    Ok(warp_protobuf::reply::protobuf(&res))
}

async fn setup_sandbox(files: Files) -> Result<TempDir, io::Error> {
    // Directory to host project:
    let tmp_dir = TempDir::new("dependagot")?;
    debug!("created: {}", tmp_dir.path().to_str().unwrap());

    // Mock src/lib.rs to be a "valid" project:
    let src_dir = tmp_dir.path().join("src");
    create_dir(&src_dir)?;
    File::create(&src_dir.join("lib.rs"))?;

    // TODO: if files contains relative paths

    // Write out files:
    let files = files.lock().await;
    for (name, data) in files.iter() {
        let file_path = tmp_dir.path().join(name);
        let mut tmp_file = File::create(&file_path)?;
        tmp_file.write_all(data.as_bytes())?;
        debug!("created: {}", file_path.to_str().unwrap());
    }
    Ok(tmp_dir)
}

fn read_new_files(sandbox: TempDir) -> Result<HashMap<String, String>, io::Error> {
    let mut new_files = HashMap::new();
    new_files.insert(
        "Cargo.toml".to_string(),
        read_to_string(sandbox.path().join("Cargo.toml"))?.to_string(),
    );
    new_files.insert(
        "Cargo.lock".to_string(),
        read_to_string(sandbox.path().join("Cargo.lock"))?.to_string(),
    );
    Ok(new_files)
}
