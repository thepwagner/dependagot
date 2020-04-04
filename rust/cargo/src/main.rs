extern crate futures;
extern crate hyper;
extern crate dependagot_common;

mod service;

use futures::Future;
use futures::sync::oneshot;
use hyper::server::Http;
use std::thread;



fn main() {
  let (shutdown_send, shutdown_recv) = oneshot::channel();

  let thread_res = thread::spawn(|| {
    println!("Starting server");
    // TODO: port from env
    let addr = "0.0.0.0:9999".parse().unwrap();
    let server = Http::new().bind(&addr,
        move || Ok(dependagot_common::UpdateService::new_server(service::UpdateService))).unwrap();
    server.run_until(shutdown_recv.map_err(|_| ())).unwrap();
    println!("Server stopped");
  });

    if let Err(err) = thread_res.join() { println!("Server panicked: {:?}", err); }
}