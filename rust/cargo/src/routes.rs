use crate::handlers;
use crate::state::State;
use std::convert::Infallible;
use warp::{Filter, Rejection, Reply};

pub fn new(state: State) -> impl Filter<Extract = impl Reply, Error = Rejection> + Clone {
    warp::path::end()
        .map(handlers::index)
        .or(
            warp::path!("twirp" / "dependabot.v1.UpdateService" / "Files")
                .and(warp::post())
                .and(warp_protobuf::body::protobuf())
                .and(with_state(state.clone()))
                .and_then(handlers::files),
        )
        .or(
            warp::path!("twirp" / "dependabot.v1.UpdateService" / "ListDependencies")
                .and(warp::post())
                .and(warp_protobuf::body::protobuf())
                .and(with_state(state.clone()))
                .and_then(handlers::list_dependencies),
        )
        .or(
            warp::path!("twirp" / "dependabot.v1.UpdateService" / "UpdateDependencies")
                .and(warp::post())
                .and(warp_protobuf::body::protobuf())
                .and(with_state(state))
                .and_then(handlers::update_dependencies),
        )
}

fn with_state(state: State) -> impl Filter<Extract = (State,), Error = Infallible> + Clone {
    warp::any().map(move || state.clone())
}
