mod files;
mod update_dependencies;

pub use files::files;
pub use update_dependencies::update_dependencies;

pub fn index() -> &'static str {
    "dependagot"
}
