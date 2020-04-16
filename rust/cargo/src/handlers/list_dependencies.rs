use crate::handlers::sandbox::setup_sandbox;
use crate::modules::state::Files;
use cargo::core::Workspace;
use cargo::util::Config;

/// ListDependencies()
pub async fn list_dependencies(
    _req: dependagot_common::ListDependenciesRequest,
    files: Files,
) -> Result<impl warp::Reply, warp::Rejection> {
    let (sandbox, _) = match setup_sandbox(files, vec![]).await {
        Err(e) => {
            error!("error creating sandbox: {:?}", e);
            // TODO: custom error
            return Err(warp::reject::not_found());
        }
        Ok(s) => s,
    };
    info!("completed sandbox: {}", sandbox.path().to_str().unwrap());

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

    let dependencies: Vec<dependagot_common::Dependency> = ws
        .current()
        .unwrap()
        .dependencies()
        .iter()
        .map(|dep| dependagot_common::Dependency {
            package: dep.package_name().to_string(),
            version: dep.version_req().to_string(),
        })
        .collect();

    let res = dependagot_common::ListDependenciesResponse { dependencies };
    Ok(warp_protobuf::reply::protobuf(&res))
}
