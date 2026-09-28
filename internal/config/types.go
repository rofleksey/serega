package config

import (
	"net"
	"time"
)

// Config is the validated process configuration used by the composition root.
type Config struct {
	HTTPAddr        string
	DatabaseURL     string
	TrustedProxies  []*net.IPNet
	ShutdownTimeout time.Duration
	LogFormat       string
	SecureCookies   bool
}
