extern crate futures;
extern crate dependagot_common;
extern crate prost_twirp;

use futures::future;

pub struct UpdateService;
impl dependagot_common::UpdateService for UpdateService {
  fn files(&self, i: dependagot_common::PTReq<dependagot_common::FilesRequest>) -> prost_twirp::PTRes<dependagot_common::FilesResponse> {
    Box::new(future::ok(
      dependagot_common::FilesResponse {
        required_paths: vec![],
        optional_paths: vec![],
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

    // fn make_hat(&self, i: service::PTReq<service::Size>) -> service::PTRes<service::Hat> {
    //     Box::new(future::ok(
    //         service::Hat { size: i.input.inches, color: "blue".to_string(), name: "fedora".to_string() }.into()
    //     ))
    // }
}