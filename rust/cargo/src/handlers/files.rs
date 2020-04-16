use crate::modules::state::Files;
use dependagot_common;
use std::str;

/// Files()
pub async fn files(
    req: dependagot_common::FilesRequest,
    files: Files,
) -> Result<impl warp::Reply, std::convert::Infallible> {
    let mut files = files.lock().await;
    for (k, v) in req.files.iter() {
        if k.as_str() == "Cargo.toml" || k.as_str() == "Cargo.lock" {
            files.insert(k.to_string(), str::from_utf8(&v).unwrap().to_string());
        }
    }

    // TODO: if Cargo.toml has been uploaded, parse for relative paths

    let mut required_paths = vec![];
    if !files.contains_key("Cargo.toml") {
        required_paths.push("Cargo.toml".to_string());
    }

    let mut optional_paths = vec![];
    if !files.contains_key("Cargo.lock") {
        optional_paths.push("Cargo.lock".to_string());
    }
    let res = dependagot_common::FilesResponse {
        required_paths,
        optional_paths,
    };
    Ok(warp_protobuf::reply::protobuf(&res))
}
