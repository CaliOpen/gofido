package config

import (
	"github.com/gocql/gocql"
)

type CassandraConfig struct {
	Hosts    []string
	Keyspace string
	// TOFIX use config value
	Consistency gocql.Consistency
}

type ServerTlsConfig struct {
	Enable   bool
	CertFile string `yaml:"cert_file"`
	KeyFile  string `yaml:"key_file"`
}

type StaticConfig struct {
	Enable    bool
	Directory string
}

type ServerConfig struct {
	ListenInterface       string          `yaml:"listen_interface"`
	ListenPort            int             `yaml:"listen_port"`
	AppId                 string          `yaml:"app_id"`
	ChallengeTtl          string          `yaml:"challenge_ttl"`
	SkipAttestationVerify string          `yaml:"skip_attestation_verify"`
	Tls                   ServerTlsConfig `yaml:"tls"`
	Static                StaticConfig    `yaml:"static"`
}

type Config struct {
	Store  CassandraConfig
	Server ServerConfig
}
