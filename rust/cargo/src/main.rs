extern crate serde_derive;
extern crate warp;

use std::str;
use warp::Filter;
use dependagot_common;

struct Modules {
    cargo_toml: String,
    cargo_lock: String,
}

#[tokio::main]
async fn main() {
    let modules = Modules {
        cargo_toml: "".to_string(),
        cargo_lock: "".to_string(),
    };
    let mutex = std::sync::Mutex::new(modules);
    let arc = std::sync::Arc::new(mutex);

    let files = warp::post()
        .and(warp::path!("twirp" / "dependabot.v1.UpdateService" / "Files"))
        .and(warp_protobuf::body::protobuf())
        .map(move |req: dependagot_common::FilesRequest| {
            let mut modules = arc.lock().unwrap();
            for (k, v) in req.files.iter() {
                if k == "Cargo.toml" {
                    modules.cargo_toml = str::from_utf8(&v).unwrap().to_string();
                } else if k == "Cargo.lock" {
                    modules.cargo_lock = str::from_utf8(&v).unwrap().to_string();
                }
            }

            let mut required_paths = vec![];
            if modules.cargo_toml == "" {
                required_paths.push("Cargo.toml".to_string());
            }

            let mut optional_paths = vec![];
            if modules.cargo_lock == "" {
                optional_paths.push("Cargo.lock".to_string());
            }
            let res = dependagot_common::FilesResponse {
                required_paths,
                optional_paths,
            };
            warp_protobuf::reply::protobuf(&res)
        });

    let list = warp::post()
        .and(warp::path!("twirp" / "dependabot.v1.UpdateService" / "ListDependencies"))
        .and(warp_protobuf::body::protobuf())
        .map(move |req: dependagot_common::ListDependenciesRequest| {
            let res = dependagot_common::ListDependenciesResponse {
                dependencies: vec![],
            };
            warp_protobuf::reply::protobuf(&res)
        });
    warp::serve(files.or(list))
        .run(([0, 0, 0, 0], 9999))
        .await;
}