package models

import "net/url"

type LoadBalancer interface {
	NextTarget() (*url.URL, error)
}
