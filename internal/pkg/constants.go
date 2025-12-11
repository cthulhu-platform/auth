package pkg

import (
	"github.com/cthulhu-platform/common/pkg/env"
)

var (
	// AMQP config
	AMQP_USER  = env.GetEnv("AMQP_USER", "guest")
	AMQP_PASS  = env.GetEnv("AMQP_PASS", "guest")
	AMQP_HOST  = env.GetEnv("AMQP_HOST", "localhost")
	AMQP_PORT  = env.GetEnv("AMQP_PORT", "5672")
	AMQP_VHOST = env.GetEnv("AMQP_VHOST", "/")

	// JWT Config
	JWT_SECRET         = env.GetEnv("JWT_SECRET", "")
	JWT_ACCESS_EXPIRY  = env.GetEnv("JWT_ACCESS_EXPIRY", "15m")
	JWT_REFRESH_EXPIRY = env.GetEnv("JWT_REFRESH_EXPIRY", "168h") // 7 days
)

