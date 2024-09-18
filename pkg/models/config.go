package models

type Config struct {
	Global GlobalConfig `hcl:"global"`
	Server ServerConfig `hcl:"server"`
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
	Path       string   `hcl:"path"`
	Handler    string   `hcl:"handler"`
	Script     string   `hcl:"script"`
	Template   string   `hcl:"template"`
	Middleware []string `hcl:"middleware"`
}

type ErrorPageConfig struct {
	StatusCode int    `hcl:"status_code"`
	Path       string `hcl:"path"`
}