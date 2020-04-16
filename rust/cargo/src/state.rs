use std::collections::HashMap;
use std::sync::Arc;
use tokio::sync::Mutex;

/// State is maintained across handler invocations.
#[derive(Clone)]
pub struct State {
    pub files: Files,
}

/// Files are paths loaded into context for dependency updates.
pub type Files = Arc<Mutex<HashMap<String, String>>>;

impl State {
    pub fn new() -> State {
        State {
            files: Arc::new(Mutex::new(HashMap::new())),
        }
    }
}
