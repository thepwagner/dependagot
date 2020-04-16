use crate::handlers::sandbox::{read_sandbox, setup_sandbox};
use crate::state::State;
use cargo::core::Workspace;
use cargo::util::Config;

pub async fn update_dependencies(
    req: dependagot_common::UpdateDependenciesRequest,
    state: State,
) -> Result<impl warp::Reply, warp::Rejection> {
    // Write files out to a temporary directory:
    let (sandbox, old_versions) = match setup_sandbox(state, req.dependencies).await {
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

    // Request cargo update the Cargo.lock file:
    match cargo::ops::update_lockfile(
        &ws,
        &cargo::ops::UpdateOptions {
            config: &config,
            aggressive: false,
            dry_run: false,
            precise: None,
            to_update: old_versions,
        },
    ) {
        Err(e) => {
            error!("error updating dependencies: {:?}", e);
            // TODO: custom error
            return Err(warp::reject::not_found());
        }
        Ok(s) => s,
    };

    // Read and return the updated files from sandbox:
    let new_files = match read_sandbox(sandbox) {
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
