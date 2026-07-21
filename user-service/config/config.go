package config

import "github.com/AnnaKhairetdinova/user-service/pkg"

type Config struct {
	DBUrl    string
	GRPCPort string
}

func Load() Config {
	cfg := Config{}

	cfg.DBUrl = pkg.MustEnv("DATABASE_URL")
	cfg.GRPCPort = pkg.EnvOr("GRPC_PORT", ":50051")

	return cfg
}
