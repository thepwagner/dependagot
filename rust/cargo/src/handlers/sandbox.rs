use crate::state::State;
use std::collections::HashMap;
use std::fs::{create_dir, read_to_string, File};
use std::io::{self, Write};
use tempdir::TempDir;
use toml::Value;

pub async fn setup_sandbox(
    state: State,
    dependencies: Vec<dependagot_common::Dependency>,
) -> Result<(TempDir, Vec<String>), io::Error> {
    // Directory to host project:
    let tmp_dir = TempDir::new("dependagot")?;
    debug!("created: {}", tmp_dir.path().to_str().unwrap());

    // Mock src/lib.rs to be a "valid" project:
    let src_dir = tmp_dir.path().join("src");
    create_dir(&src_dir)?;
    File::create(&src_dir.join("lib.rs"))?;

    // TODO: if files contains relative paths

    // Write out files, excluding Cargo.toml:
    let files = state.files.lock().await;
    for (name, data) in files.iter() {
        if name == "Cargo.toml" {
            continue;
        }
        let file_path = tmp_dir.path().join(name);
        let mut tmp_file = File::create(&file_path)?;
        tmp_file.write_all(data.as_bytes())?;
        debug!("created: {}", file_path.to_str().unwrap());
    }

    // Index upgrade targets:
    let targets: HashMap<String, String> = dependencies
        .into_iter()
        .map(|dep| (dep.package, dep.version))
        .collect();

    // Iterate the [dependencies] section, replacing any requesting updates:
    let mut new_dependencies = toml::map::Map::new();
    let mut old_versions = vec![];
    let mut cargo_toml = files.get("Cargo.toml").unwrap().parse::<Value>().unwrap();
    let dependencies = cargo_toml["dependencies"].as_table().unwrap();
    for (dep, req) in dependencies.iter() {
        let value: Value = match targets.get(dep) {
            Some(target) => {
                let old_version = req.as_str().unwrap();
                info!("updating dependency {}: {} -> {}", dep, old_version, target);
                old_versions.push(format!("{}:{}", dep, old_version));
                Value::String(target.to_string())
            }
            None => Value::String(req.as_str().unwrap().to_string()),
        };
        new_dependencies.insert(dep.to_string(), value);
    }
    cargo_toml["dependencies"] = Value::Table(new_dependencies);

    // Write the edited Cargo.toml:
    let cargo_toml_path = tmp_dir.path().join("Cargo.toml");
    let mut cargo_toml_out = File::create(&cargo_toml_path)?;
    cargo_toml_out.write_all(toml::to_string(&cargo_toml).unwrap().as_bytes())?;
    debug!("created: {}", cargo_toml_path.to_str().unwrap());

    Ok((tmp_dir, old_versions))
}

pub fn read_sandbox(sandbox: TempDir) -> Result<HashMap<String, String>, io::Error> {
    let mut new_files = HashMap::new();
    new_files.insert(
        "Cargo.toml".to_string(),
        read_to_string(sandbox.path().join("Cargo.toml"))?.to_string(),
    );
    new_files.insert(
        "Cargo.lock".to_string(),
        read_to_string(sandbox.path().join("Cargo.lock"))?.to_string(),
    );
    Ok(new_files)
}
