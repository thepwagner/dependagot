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
                warp::path!("twirp" / "dependabot.v1.UpdateService" / "ListDependencies")
                    .and(warp::post())
                    .and(warp_protobuf::body::protobuf())
                    .and(with_files(state.clone()))
                    .and_then(handlers::list_dependencies),
            )
            .or(
                warp::path!("twirp" / "dependabot.v1.UpdateService" / "UpdateDependencies")
                    .and(warp::post())
                    .and(warp_protobuf::body::protobuf())
                    .and(with_files(state))
                    .and_then(handlers::update_dependencies),
            )
    }

    fn with_files(files: Files) -> impl Filter<Extract = (Files,), Error = Infallible> + Clone {
        warp::any().map(move || files.clone())
    }
}
