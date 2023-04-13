package router

import (
	"os"

	"github.com/spf13/pflag"

	"zq-xu/helper/config"
	"zq-xu/helper/log"
)

const (
	routerConfigName = "RouterConfig"

	TLSKeyPathEnv  = "TLSKeyPath"
	TLSCertPathEnv = "TLSCertPath"
)

type RouteConfig struct {
	IP        string
	Port      string
	PprofPort string

	DisableTLS bool

	// For develop, use the command below to generate the private key and cert:
	//     for key:  openssl genrsa -out server.key 2048
	//     for cert: openssl req -new -x509 -key server.key -out server.pem -days 3650
	KeyPath  string
	CertPath string
}

var (
	RouteCfg = &RouteConfig{}
)

func init() {
	config.RegisterCfg(routerConfigName, RouteCfg)
}

// AddFlags adds flags for router
func (rc *RouteConfig) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&rc.IP, "router-ip", "0.0.0.0", "the ip for listening")
	fs.StringVar(&rc.Port, "router-port", "8080", "the port for listening")
	fs.StringVar(&rc.PprofPort, "pprof-port", "6069", "the port for pprof")

	fs.BoolVar(&rc.DisableTLS, "disable-tls", false, "disable the TLS")
	fs.StringVar(&rc.KeyPath, "key-path", os.Getenv(TLSKeyPathEnv), "the key file pah of the tls")
	fs.StringVar(&rc.CertPath, "cert-path", os.Getenv(TLSCertPathEnv), "the cert file pah of the tls")
}

func (rc *RouteConfig) Revise() {
	log.Logger.Infof("RouterConfig is IP: %v; Port: %v; PprofPort: %v; DisableTLS: %v",
		rc.IP, rc.Port, rc.PprofPort, rc.DisableTLS)
}
