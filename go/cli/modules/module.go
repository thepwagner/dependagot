package modules

type DependencyModule string

const (
	GoModules   DependencyModule = "go:modules"
	RubyBundler DependencyModule = "ruby:bundler"
)

var Modules = []DependencyModule{
	GoModules,
	RubyBundler,
}
