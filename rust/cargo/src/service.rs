extern crate futures;
extern crate dependagot_common;
extern crate prost_twirp;

use futures::future;
use std::str;

#[derive(Clone)]
pub struct UpdateService {
  pub cargo_toml: String,
  pub cargo_lock: String,
}

impl dependagot_common::UpdateService for UpdateService {
  fn files(&mut self, i: dependagot_common::PTReq<dependagot_common::FilesRequest>) -> prost_twirp::PTRes<dependagot_common::FilesResponse> {
    println!("{:#?}", i);
    for (path, data) in i.input.files.iter() {
      match path.as_str() {
        "Cargo.toml" => self.cargo_toml = str::from_utf8(data).unwrap().to_string(),
        "Cargo.lock" => self.cargo_lock = str::from_utf8(data).unwrap().to_string()
      }
    }

    let mut required_paths = vec![];
    let mut optional_paths = vec![];
    if self.cargo_toml.is_empty() {
      required_paths.push("Cargo.toml".to_string())
    }
    if self.cargo_lock.is_empty() {
      required_paths.push("Cargo.lock".to_string())
    }
    Box::new(future::ok(
      dependagot_common::FilesResponse {
        required_paths: required_paths,
        optional_paths: optional_paths,
      }.into()
    ))
  }

  fn list_dependencies(&self, i: dependagot_common::PTReq<dependagot_common::ListDependenciesRequest>) -> dependagot_common::PTRes<dependagot_common::ListDependenciesResponse> {
    Box::new(future::ok(
      dependagot_common::ListDependenciesResponse {
        dependencies: vec![],
      }.into()
    ))
  }

  fn update_dependencies(&self, i: dependagot_common::PTReq<dependagot_common::UpdateDependenciesRequest>) -> dependagot_common::PTRes<dependagot_common::UpdateDependenciesResponse> {
    Box::new(future::ok(
      dependagot_common::UpdateDependenciesResponse {
        new_files: ::std::collections::HashMap::new(),
      }.into()
    ))
  }
}