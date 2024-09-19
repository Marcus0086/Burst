package models

type Config struct {
	Global  GlobalConfig `hcl:"global"`
	Server  ServerConfig `hcl:"server"`
	Include []string     `hcl:"include"`
}

type GlobalConfig struct {
	Workers int `hcl:"workers"`
}

type ServerConfig struct {
	Listen     string            `hcl:"listen"`
	HTTPS      bool              `hcl:"https"`
	Root       string            `hcl:"root"`
	Routes     []RouteConfig     `hcl:"routes"`
	ErrorPages []ErrorPageConfig `hcl:"error_pages"`
	Middleware []string          `hcl:"middleware"`
}

type RouteConfig struct {
	Path          string                 `hcl:"path"`
	Handler       string                 `hcl:"handler"`
	Template      string                 `hcl:"template,optional"`
	Body          string                 `hcl:"body,optional"`
	Middleware    []string               `hcl:"middleware,optional"`
	Data          map[string]interface{} `hcl:"data,optional"`
	Target        string                 `hcl:"target,optional"`
	Targets       []string               `hcl:"targets,optional"`
	LoadBalancing string                 `hcl:"load_balancing,optional"`
}

type ErrorPageConfig struct {
	StatusCode int    `hcl:"status_code"`
	Path       string `hcl:"path"`
}
