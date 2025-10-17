package routerx

var RouterConfig Config

type Config struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`

	DisableTLS bool `yaml:"disableTLS"`

	// For develop, use the command below to generate the private key and cert:
	//     for key:  openssl genrsa -out server.key 2048
	//     for cert: openssl req -new -x509 -key server.key -out server.pem -days 3650
	KeyPath  string `yaml:"keyPath"`
	CertPath string `yaml:"certPath"`
}

func (rc *Config) OptDefalt() {
	if rc.Host == "" {
		rc.Host = "0.0.0.0"
	}

	if rc.Port == 0 {
		rc.Port = 8080
	}
}
