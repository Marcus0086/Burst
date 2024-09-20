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
	Handler       Handler                `hcl:"handler"`
	Template      string                 `hcl:"template,optional"`
	Body          string                 `hcl:"body,optional"`
	Middleware    []string               `hcl:"middleware,optional"`
	Data          map[string]interface{} `hcl:"data,optional"`
	Target        string                 `hcl:"target,optional"`
	Targets       []string               `hcl:"targets,optional"`
	LoadBalancing string                 `hcl:"load_balancing,optional"`
	ReqHeaders    map[string]interface{} `hcl:"req_headers,optional"`
	RespHeaders   map[string]interface{} `hcl:"resp_headers,optional"`
	Redirects     []RedirectConfig       `hcl:"redirects,optional"`
}

type Handler string

const (
	StaticHandler Handler = "static"
	StaticContentHandler Handler = "static_content"
	DynamicHandler Handler = "dynamic"
	ReverseProxyHandler Handler = "reverse_proxy"
)

type RedirectConfig struct {
	From    string `hcl:"from"`
	To      string `hcl:"to"`
	Status  int    `hcl:"status"`
	Message string `hcl:"message"`
}

type ErrorPageConfig struct {
	StatusCode int    `hcl:"status_code"`
	Path       string `hcl:"path"`
}
