extern crate futures;
extern crate hyper;
extern crate prost;
#[macro_use]
extern crate prost_derive;

mod service {
  include!(concat!(env!("OUT_DIR"), "/dependagot.v1.rs"));
}
pub use service::*;
