package config

import (
	"strings"
	"testing"
)

func TestAppConfigValidate_RequiredValues(t *testing.T) {
	valid := AppConfig{
		Server:    ServerConfig{Port: "8082"},
		MongoDB:   MongoDBConfig{DSN: "mongodb://example.invalid", Database: "airis"},
		SecretKey: "secret",
	}

	tests := []struct {
		name    string
		mutate  func(*AppConfig)
		wantErr string
	}{
		{name: "valid"},
		{name: "missing port", mutate: func(config *AppConfig) { config.Server.Port = "" }, wantErr: "server port"},
		{name: "missing Mongo DSN", mutate: func(config *AppConfig) { config.MongoDB.DSN = "" }, wantErr: "MongoDB DSN"},
		{name: "missing Mongo database", mutate: func(config *AppConfig) { config.MongoDB.Database = "" }, wantErr: "MongoDB database"},
		{name: "missing secret", mutate: func(config *AppConfig) { config.SecretKey = "" }, wantErr: "secret key"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			candidate := valid
			if test.mutate != nil {
				test.mutate(&candidate)
			}
			err := candidate.Validate()
			if test.wantErr == "" {
				if err != nil {
					t.Fatalf("Validate() error = %v", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), test.wantErr) {
				t.Fatalf("Validate() error = %v, want text %q", err, test.wantErr)
			}
		})
	}
}

func TestLoad_ReadsEnvironmentAndDefaults(t *testing.T) {
	previous := Config
	t.Cleanup(func() { Config = previous })

	t.Setenv("SERVER_PORT", "")
	t.Setenv("ENV", "test")
	t.Setenv("MONGODB_DSN", "mongodb://example.invalid")
	t.Setenv("MONGODB_DATABASE", "airis_test")
	t.Setenv("MONGODB_COLLECTION", "")
	t.Setenv("MONGODB_MAX_POOL", "invalid")
	t.Setenv("MONGODB_MIN_POOL", "3")
	t.Setenv("SECRET_KEY", "test-secret")

	if err := Load(); err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if Config.Server.Port != "8082" || Config.Server.Env != "test" {
		t.Fatalf("server config = %#v", Config.Server)
	}
	if Config.MongoDB.Collection != "data_20251101_0" {
		t.Fatalf("collection = %q", Config.MongoDB.Collection)
	}
	if Config.MongoDB.MaxPool != 100 || Config.MongoDB.MinPool != 3 {
		t.Fatalf("pool config = max %d, min %d", Config.MongoDB.MaxPool, Config.MongoDB.MinPool)
	}
	if err := Config.Validate(); err != nil {
		t.Fatalf("loaded config Validate() error = %v", err)
	}
}
