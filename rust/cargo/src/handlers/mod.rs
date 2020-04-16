mod files;
mod list_dependencies;
mod update_dependencies;
mod sandbox;

pub use files::files;
pub use list_dependencies::list_dependencies;
pub use update_dependencies::update_dependencies;

pub fn index() -> &'static str {
    "dependagot"
}
