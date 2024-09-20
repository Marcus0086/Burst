# Burst

Burst is a super minimal, easy-to-use server and reverse proxy with modern features, designed for simplicity and performance.

## Features

- Lightweight and fast
- Simple configuration using Burstfile
- Reverse proxy capabilities
- Modern features for efficient serving
- Memory and thread-safe
- Optimized with Go concurrency

## Quick Start

1. Build Burst

```
git clone https://github.com/Marcus0086/Burst
cd Burst
# Create a Burstfile in the root of your project.
# To add more files, create a config directory and add *.burst files in it.
# Build Burst
go build -ldflags "-s -w" -o burst
```

2. Create a `Burstfile` in your project directory
3. Run `burst` command in the terminal

## Burstfile Example

```
workers = 10    

server {
    root = "./public"
    listen = ":80"
    routes = [
        {
            path = "/*"
            handler = "static"
        }
    ]
}
```

## More Features

Burst supports static file serving, dynamic content rendering, and reverse proxy capabilities.

### Static File Serving

```
server {
    root = "./public"
    listen = ":80"
    routes = [
        {
            path = "/*"
            handler = "static"
        }
    ]
}
```

### Dynamic Content Rendering

```
server {
    root = "./public"
    listen = ":80"
    routes = [
        {
            path "/data"
            handler = "dynamic"
            template = "data.html"
            data = {
                "name" = "John"
                "age" = 30
            }
        }
    ]
}
```

### Reverse Proxy

```
server {
    listen = ":80"
    routes = [
        {
            path = "/*"
            handler = "reverse_proxy"
            target = "http://localhost:3000"
        }
    ]
}
```

### Load Balancing

```
server {
    listen = ":80"
    routes = [
        {
            path = "/*"
            handler = "reverse_proxy"
            targets = ["http://localhost:3000", "http://localhost:3001"]
            load_balancing = "round_robin"
        }
    ]
}
```

### Redirects

```
server {
    listen = ":80"
    routes = [
        {
            path = "/old"
            redirects = [
                {
                    from = "/old/*"
                    to = "/new/*"
                    status = 301
                    message = "Moved Permanently"
                }
            ]   
        }
    ]
}
```

### Request Headers

```
server {
    listen = ":80"
    routes = [
        {
            path = "/*"
            handler = "reverse_proxy"
            target = "http://localhost:3000"
            req_headers = {
                "Authorization" = "Bearer token"
                "User-Agent" = false (disable)
            }
        }
    ]
}
```

### Response Headers

```
server {
    listen = ":80"
    routes = [
        {
            path = "/*"
            handler = "reverse_proxy"
            target = "http://localhost:3000"
            res_headers = {
                "Content-Type" = "text/html; charset=utf-8"
                "Server" = false (disable)
            }
        }
    ]
}
```

