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
    use super::handlers;
    use super::state::Files;
    use std::convert::Infallible;
    use warp::{Filter, Rejection, Reply};

    pub fn new(state: Files) -> impl Filter<Extract = impl Reply, Error = Rejection> + Clone {
        files(state.clone())
            .or(list_dependencies(state.clone()))
            .or(update_dependencies(state))
    }

    fn files(state: Files) -> impl Filter<Extract = impl Reply, Error = Rejection> + Clone {
        warp::path!("twirp" / "dependabot.v1.UpdateService" / "Files")
            .and(warp::post())
            .and(warp_protobuf::body::protobuf())
            .and(with_files(state))
            .and_then(handlers::files)
    }

    fn list_dependencies(
        state: Files,
    ) -> impl Filter<Extract = impl Reply, Error = Rejection> + Clone {
        warp::path!("twirp" / "dependabot.v1.UpdateService" / "ListDependencies")
            .and(warp::post())
            .and(warp_protobuf::body::protobuf())
            .and(with_files(state))
            .and_then(handlers::list_dependencies)
    }

    fn update_dependencies(
        state: Files,
    ) -> impl Filter<Extract = impl Reply, Error = Rejection> + Clone {
        warp::path!("twirp" / "dependabot.v1.UpdateService" / "UpdateDependencies")
            .and(warp::post())
            .and(warp_protobuf::body::protobuf())
            .and(with_files(state))
            .and_then(handlers::update_dependencies)
    }

    fn with_files(files: Files) -> impl Filter<Extract = (Files,), Error = Infallible> + Clone {
        warp::any().map(move || files.clone())
    }
}

mod handlers {
    extern crate cargo_toml;

    use super::state::Files;
    use cargo_toml::Manifest;
    use dependagot_common;
    use std::collections::HashMap;
    use std::convert::Infallible;
    use std::str;
    use warp::Reply;

    pub async fn files(
        req: dependagot_common::FilesRequest,
        files: Files,
    ) -> Result<impl Reply, Infallible> {
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

    pub async fn update_dependencies(
        _req: dependagot_common::UpdateDependenciesRequest,
        files: Files,
    ) -> Result<impl warp::Reply, warp::Rejection> {
        let files = files.lock().await;

        let new_files = HashMap::new();

        let res = dependagot_common::UpdateDependenciesResponse { new_files };
        Ok(warp_protobuf::reply::protobuf(&res))
    }
}
