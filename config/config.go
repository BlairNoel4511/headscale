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
	// EphemeralNodeInactivityTimeout controls how long an ephemeral node can be
	// inactive before it is removed. Increased default to 360s for my home lab
	// setup where nodes occasionally go to sleep for longer periods.
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
	ExtraParams