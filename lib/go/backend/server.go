package backend

import (
	libconfig "lib/config"
	"lib/httpserver"
)

// ServerFrom assembles the server settings from the env blocks every
// Backend embeds (lib/config.Server and CORS) and its pprof port.
func ServerFrom(s libconfig.Server, cors libconfig.CORS, profilePort string) Server {
	return Server{
		Port:               s.Port,
		ProfilePort:        profilePort,
		CorsAllowedOrigins: cors.AllowedOrigins,
		CorsAllowedHeaders: cors.AllowedHeaders,
		RateLimit: httpserver.RateLimit{
			Disabled: s.DisableRateLimit,
			RPS:      s.RateLimitRPS,
			Burst:    s.RateLimitBurst,
		},
	}
}
