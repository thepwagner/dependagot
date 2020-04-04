package modules

type DependencyModule string

const (
	GoModules   DependencyModule = "go:modules"
	RubyBundler DependencyModule = "ruby:bundler"
	RustCargo DependencyModule = "rust:cargo"
)

var Modules = []DependencyModule{
	GoModules,
	RubyBundler,
	RustCargo,
}
