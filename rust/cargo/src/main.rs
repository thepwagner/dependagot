// extern crate serde_derive;
// extern crate warp;
#[macro_use]
extern crate log;

use std::env;
use warp::Filter;

mod modules;

#[tokio::main]
async fn main() {
    if env::var_os("RUST_LOG").is_none() {
        env::set_var("RUST_LOG", "dependagot=info");
    }
    env_logger::init();

    let files = modules::state::empty_files();
    let api = modules::filters::new(files);

    // TODO: from env
    let port = 9999;
    let routes = api.with(warp::log("dependagot"));
    info!("starting server 0.0.0.0:{}", port);
    warp::serve(routes).run(([0, 0, 0, 0], port)).await;
}
