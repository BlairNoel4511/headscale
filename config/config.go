// Package config provides configuration management for headscale.
package config

import (
	"fmt"
	"net/netip"
	"os"
	"time"

	"github.com/spf13/viper"
)

// Config holds the full headscale server configuration.
type Config struct {
	ServerURL          string        `mapstructure:"server_url"`
	ListenAddr         string        `mapstructure:"listen_addr"`
	MetricsListenAddr  string        `mapstructure:"metrics_listen_addr"`
	GRPCListenAddr     string        `mapstructure:"grpc_listen_addr"`
	GRPCAllowInsecure  bool          `mapstructure:"grpc_allow_insecure"`
	PrivateKeyPath     string        `mapstructure:"private_key_path"`
	NoisePrivateKeyPath string       `mapstructure:"noise_private_key_path"`
	DBtype             string        `mapstructure:"db_type"`
	DBpath             string        `mapstructure:"db_path"`
	DBhost             string        `mapstructure:"db_host"`
	DBport             int           `mapstructure:"db_port"`
	DBname             string        `mapstructure:"db_name"`
	DBuser             string        `mapstructure:"db_user"`
	DBpass             string        `mapstructure:"db_pass"`
	DBssl              string        `mapstructure:"db_ssl"`
	TLSLetsEncryptHostname  string   `mapstructure:"tls_letsencrypt_hostname"`
	TLSLetsEncryptCacheDir  string   `mapstructure:"tls_letsencrypt_cache_dir"`
	TLSLetsEncryptChallenge string   `mapstructure:"tls_letsencrypt_challenge_type"`
	TLSCertPath        string        `mapstructure:"tls_cert_path"`
	TLSKeyPath         string        `mapstructure:"tls_key_path"`
	ACMEEmail          string        `mapstructure:"acme_email"`
	ACMEUrl            string        `mapstructure:"acme_url"`
	DNSConfig          *DNSConfig    `mapstructure:"dns_config"`
	IPPrefixes         []netip.Prefix `mapstructure:"ip_prefixes"`
	BaseDomain         string        `mapstructure:"base_domain"`
	LogLevel           string        `mapstructure:"log_level"`
	DisableCheckUpdates bool         `mapstructure:"disable_check_updates"`
	EphemeralNodeInactivityTimeout time.Duration `mapstructure:"ephemeral_node_inactivity_timeout"`
	NodeUpdateCheckInterval        time.Duration `mapstructure:"node_update_check_interval"`
	OIDC               OIDCConfig    `mapstructure:"oidc"`
}

// DNSConfig holds DNS-related settings.
type DNSConfig struct {
	OverrideLocalDNS bool     `mapstructure:"override_local_dns"`
	Nameservers      []string `mapstructure:"nameservers"`
	RestrictedNameservers map[string][]string `mapstructure:"restricted_nameservers"`
	Domains          []string `mapstructure:"domains"`
	MagicDNS         bool     `mapstructure:"magic_dns"`
	BaseDomain       string   `mapstructure:"base_domain"`
}

// OIDCConfig holds OpenID Connect settings.
type OIDCConfig struct {
	Issuer           string            `mapstructure:"issuer"`
	ClientID         string            `mapstructure:"client_id"`
	ClientSecret     string            `mapstructure:"client_secret"`
	Scope            []string          `mapstructure:"scope"`
	ExtraParams      map[string]string `mapstructure:"extra_params"`
	AllowedDomains   []string          `mapstructure:"allowed_domains"`
	AllowedUsers     []string          `mapstructure:"allowed_users"`
	StripEmaildomain bool              `mapstructure:"strip_email_domain"`
}

// LoadConfig reads the configuration from the given path using viper.
func LoadConfig(path string, isFile bool) error {
	if isFile {
		viper.SetConfigFile(path)
	} else {
		viper.SetConfigName("config")
		if path == "" {
			viper.AddConfigPath("/etc/headscale/")
			viper.AddConfigPath("$HOME/.headscale")
			viper.AddConfigPath(".")
		} else {
			viper.AddConfigPath(path)
		}
	}

	viper.SetEnvPrefix("headscale")
	viper.AutomaticEnv()

	setDefaults()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("reading config: %w", err)
		}
	}

	return nil
}

// GetConfig unmarshals the viper configuration into a Config struct.
func GetConfig() (*Config, error) {
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshaling config: %w", err)
	}

	return &cfg, nil
}

// setDefaults applies sane default values for all configuration keys.
func setDefaults() {
	viper.SetDefault("listen_addr", "0.0.0.0:8080")
	viper.SetDefault("metrics_listen_addr", "127.0.0.1:9090")
	viper.SetDefault("grpc_listen_addr", "0.0.0.0:50443")
	viper.SetDefault("grpc_allow_insecure", false)
	viper.SetDefault("db_type", "sqlite3")
	viper.SetDefault("db_path", "/var/lib/headscale/db.sqlite")
	viper.SetDefault("private_key_path", "/var/lib/headscale/private.key")
	viper.SetDefault("noise_private_key_path", "/var/lib/headscale/noise_private.key")
	viper.SetDefault("log_level", "info")
	viper.SetDefault("ip_prefixes", []string{"100.64.0.0/10", "fd7a:115c:a1e0::/48"})
	viper.SetDefault("ephemeral_node_inactivity_timeout", "120s")
	viper.SetDefault("node_update_check_interval", "10s")
	viper.SetDefault("disable_check_updates", false)
	viper.SetDefault("tls_letsencrypt_cache_dir", "/var/www/.cache")
	viper.SetDefault("tls_letsencrypt_challenge_type", "HTTP-01")
	viper.SetDefault("oidc.scope", []string{"openid", "profile", "email"})
	viper.SetDefault("oidc.strip_email_domain", true)

	// Expand environment variables in paths
	viper.SetDefault("db_path", os.ExpandEnv(viper.GetString("db_path")))
}
