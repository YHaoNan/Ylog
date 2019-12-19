package plugins

type Plugin interface {
	Name() string
	IsCommandLinePlugin() bool
	Exec(params map[string]interface{})
}

