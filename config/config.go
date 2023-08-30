package config

import (
	"bytes"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/grip211/order/pkg/database"
)

type Server struct {
	Database            database.Opt `yaml:"database"`
	BugsnagAPIKey       string       `yaml:"bugsnag_api_key"`
	BugsnagReleaseStage string       `yaml:"bugsnag_release_stage"`
	MaxMindDatabaseDir  string       `yaml:"max_mind_database_dir"`
	EnableV1EPG         bool         `yaml:"enable_v1_epg"`
	EnableV4EPG         bool         `yaml:"enable_v4_epg"`
	Prefork             bool         `yaml:"prefork"`
	NatsURL             string       `yaml:"nats_url"`
}

type Config struct {
	Server `yaml:"server"`
}

func New(filepath string) (*Config, error) {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	config := &Config{}
	d := yaml.NewDecoder(bytes.NewReader(content))
	if err = d.Decode(&config); err != nil {
		return nil, err
	}
	return config, nil
}
