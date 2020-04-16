// extern crate serde_derive;
// extern crate warp;
#[macro_use]
extern crate log;

use std::env;
use warp::Filter;

mod handlers;
mod routes;
mod state;

#[tokio::main]
async fn main() {
    if env::var_os("RUST_LOG").is_none() {
        env::set_var("RUST_LOG", "dependagot=info");
    }
    env_logger::init();

    let state = state::State::new();
    let port: u16 = match env::var_os("DEPENDAGOT_PORT") {
        None => 9999,
        Some(port) => port.into_string().unwrap().parse::<u16>().unwrap(),
    };
    let routes = routes::new(state).with(warp::log("dependagot"));
    info!("starting server 0.0.0.0:{}", port);
    warp::serve(routes).run(([0, 0, 0, 0], port)).await;
}
