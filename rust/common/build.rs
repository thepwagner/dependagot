extern crate prost_build;

fn main() {
    let mut conf = prost_build::Config::new();
    conf.compile_protos(&["dependabot.proto"], &["dependabot/v1/"]).unwrap();
}
