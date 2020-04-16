mod files;
mod list_dependencies;
mod sandbox;
mod update_dependencies;

pub use files::files;
pub use list_dependencies::list_dependencies;
pub use update_dependencies::update_dependencies;

pub fn index() -> &'static str {
    "dependagot"
}
