package config

// Config holds service-level configuration.
type Config struct {
	HTTPPort    string
	Environment string
}

// Default returns sensible defaults for local development and demos.
func Default() Config {
	return Config{
		HTTPPort:    ":8080",
		Environment: "development",
	}
}
