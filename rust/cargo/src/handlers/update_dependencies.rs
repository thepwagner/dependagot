use crate::modules::state::Files;
use std::collections::HashMap;
use std::fs::File;
use std::io::{self, Write};
use tempdir::TempDir;

pub async fn update_dependencies(
    _req: dependagot_common::UpdateDependenciesRequest,
    files: Files,
) -> Result<impl warp::Reply, warp::Rejection> {
    let sandbox = match setup_sandbox(files).await {
        Err(e) => {
            error!("error creating sandbox: {:?}", e);
            // TODO: custom error
            return Err(warp::reject::not_found());
        }
        Ok(s) => s,
    };

    info!("sandbox: {}", sandbox.path().to_str().unwrap());

    let mut new_files = HashMap::new();
    new_files.insert("foo".to_string(), "bar".to_string());

    let res = dependagot_common::UpdateDependenciesResponse { new_files };
    Ok(warp_protobuf::reply::protobuf(&res))
}

async fn setup_sandbox(files: Files) -> Result<TempDir, io::Error> {
    let tmp_dir = TempDir::new("dependagot")?;
    debug!("created: {}", tmp_dir.path().to_str().unwrap());

    let files = files.lock().await;
    for (name, data) in files.iter() {
        let file_path = tmp_dir.path().join(name);
        let mut tmp_file = File::create(&file_path)?;
        tmp_file.write_all(data.as_bytes())?;
        debug!("created: {}", file_path.to_str().unwrap());
    }
    Ok(tmp_dir)
}
