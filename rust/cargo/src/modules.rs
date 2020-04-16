extern crate warp;

pub mod state {
    use std::collections::HashMap;
    use std::sync::Arc;
    use tokio::sync::Mutex;

    pub type Files = Arc<Mutex<HashMap<String, String>>>;
    pub fn empty_files() -> Files {
        Arc::new(Mutex::new(HashMap::new()))
    }
}

pub mod filters {
    use super::state::Files;
    use crate::handlers;
    use std::convert::Infallible;
    use warp::{Filter, Rejection, Reply};

    pub fn new(state: Files) -> impl Filter<Extract = impl Reply, Error = Rejection> + Clone {
        warp::path::end()
            .map(handlers::index)
            .or(
                warp::path!("twirp" / "dependabot.v1.UpdateService" / "Files")
                    .and(warp::post())
                    .and(warp_protobuf::body::protobuf())
                    .and(with_files(state.clone()))
                    .and_then(handlers::files),
            )
            .or(
                warp::path!("twirp" / "dependabot.v1.UpdateService" / "UpdateDependencies")
                    .and(warp::post())
                    .and(warp_protobuf::body::protobuf())
                    .and(with_files(state))
                    .and_then(handlers::update_dependencies),
            )
        // .or(list_dependencies(state.clone()))
        // .or(update_dependencies(state))
    }

    // fn list_dependencies(
    //     state: Files,
    // ) -> impl Filter<Extract = impl Reply, Error = Rejection> + Clone {

    // }

    // fn update_dependencies(
    //     state: Files,
    // ) -> impl Filter<Extract = impl Reply, Error = Rejection> + Clone {
    //     warp::path!("twirp" / "dependabot.v1.UpdateService" / "UpdateDependencies")
    //         .and(warp::post())
    //         .and(warp_protobuf::body::protobuf())
    //         .and(with_files(state))
    //         .and_then(handlers::update_dependencies)
    // }

    fn with_files(files: Files) -> impl Filter<Extract = (Files,), Error = Infallible> + Clone {
        warp::any().map(move || files.clone())
    }
}

mod handlers {
    extern crate cargo_toml;

    use super::state::Files;
    use cargo_toml::Manifest;
    use dependagot_common;

    pub async fn list_dependencies(
        _req: dependagot_common::ListDependenciesRequest,
        files: Files,
    ) -> Result<impl warp::Reply, warp::Rejection> {
        let files = files.lock().await;
        if let Some(toml) = files.get("Cargo.toml") {
            if let Some(manifest) = Manifest::from_str(toml).ok() {
                let mut dependencies = vec![];
                for (name, dep) in manifest.dependencies.iter() {
                    dependencies.push(dependagot_common::Dependency {
                        package: name.to_string(),
                        version: dep.req().to_string(),
                    });
                }

                let res = dependagot_common::ListDependenciesResponse { dependencies };
                Ok(warp_protobuf::reply::protobuf(&res))
            } else {
                // TODO: invalid manifest error
                Err(warp::reject::not_found())
            }
        } else {
            // TODO: invalid state error
            Err(warp::reject::not_found())
        }
    }
}
