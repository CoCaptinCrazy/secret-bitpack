package config

type Config struct {
	Environment string `json:"environment"`
	Registry    string `json:"registry"`
	DataDir     string `json:"data_dir"`
}

func Default() *Config {
	return &Config{
		Environment: "development",
		Registry:    "registry.local",
		DataDir:     "data",
	}
}
