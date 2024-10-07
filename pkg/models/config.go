package models

type Config struct {
	Global  GlobalConfig `hcl:"global"`
	Server  ServerConfig `hcl:"server"`
	Include []string     `hcl:"include,optional"`
}

type ConfigJSON struct {
	Apps Apps `json:"apps"`
}

type GlobalConfig struct {
	Workers  int  `hcl:"workers"`
	AdminAPI bool `hcl:"admin_api,optional"`
}

type Apps struct {
	HTTP HTTPApp `json:"http,omitempty"`
}

type HTTPApp struct {
	Servers map[string]ServerConfig `hcl:"servers" json:"servers,omitempty"`
}

type ServerConfig struct {
	Listen     string            `hcl:"listen" json:"listen,omitempty"`
	HTTPS      bool              `hcl:"https,optional" json:"https,omitempty"`
	Root       string            `hcl:"root,optional" json:"root,omitempty"`
	Routes     []RouteConfig     `hcl:"routes" json:"routes,omitempty"`
	ErrorPages []ErrorPageConfig `hcl:"error_pages,optional" json:"error_pages,omitempty"`
	Middleware []string          `hcl:"middleware,optional" json:"middleware,omitempty"`
}

type RouteConfig struct {
	Path          string                 `hcl:"path" json:"path,omitempty"`
	Handler       Handler                `hcl:"handler" json:"handler,omitempty"`
	Template      string                 `hcl:"template,optional" json:"template,omitempty"`
	Body          string                 `hcl:"body,optional" json:"body,omitempty"`
	Middleware    []string               `hcl:"middleware,optional" json:"middleware,omitempty"`
	Data          map[string]interface{} `hcl:"data,optional" json:"data,omitempty"`
	Target        string                 `hcl:"target,optional" json:"target,omitempty"`
	Targets       []string               `hcl:"targets,optional" json:"targets,omitempty"`
	LoadBalancing string                 `hcl:"load_balancing,optional" json:"load_balancing,omitempty"`
	ReqHeaders    map[string]interface{} `hcl:"req_headers,optional" json:"req_headers,omitempty"`
	RespHeaders   map[string]interface{} `hcl:"resp_headers,optional" json:"resp_headers,omitempty"`
	Redirects     []RedirectConfig       `hcl:"redirects,optional" json:"redirects,omitempty"`
	Method        string                 `hcl:"method,optional" json:"method,omitempty"`
}

type Handler string

const (
	StaticHandler        Handler = "static"
	StaticContentHandler Handler = "static_content"
	DynamicHandler       Handler = "dynamic"
	ReverseProxyHandler  Handler = "reverse_proxy"
	AdminAPIHandler      Handler = "admin_api"
)

type RedirectConfig struct {
	From    string `hcl:"from" json:"from,omitempty"`
	To      string `hcl:"to" json:"to,omitempty"`
	Status  int    `hcl:"status" json:"status,omitempty"`
	Message string `hcl:"message" json:"message,omitempty"`
}

type ErrorPageConfig struct {
	StatusCode int    `hcl:"status_code" json:"status_code,omitempty"`
	Path       string `hcl:"path" json:"path,omitempty"`
}
