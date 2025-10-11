package awsapi

import (
	"os"

	"github.com/spf13/pflag"

	"zq-xu/y-travel-go/pkg/config"
	"zq-xu/y-travel-go/pkg/logs"
)

const (
	awsSessionConfigName = "AWSSessionConfig"

	AWSAccessIDEnv     = "AWSAccessID"
	AWSAccessSecretEnv = "AWSAccessSecret"
	AWSRegionEnv       = "AWSRegion"
	AWSEndpointEnv     = "AWSEndpoint"
)

var (
	SessionCfg = &SessionConfig{}
)

type SessionConfig struct {
	AccessID     string
	AccessSecret string
	Region       string
	Endpoint     string
}

func InitSessionConfig() {
	config.RegisterCfg(awsSessionConfigName, SessionCfg)
}

// AddFlags adds flags for SessionConfig
func (sc *SessionConfig) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&sc.AccessID, "aws-access-id", os.Getenv(AWSAccessIDEnv), "the AccessID.")
	fs.StringVar(&sc.AccessSecret, "aws-access-secret", os.Getenv(AWSAccessSecretEnv), "the AccessSecret.")
	fs.StringVar(&sc.Region, "aws-region", os.Getenv(AWSRegionEnv), "the Region.")
	fs.StringVar(&sc.Endpoint, "aws-endpoint", os.Getenv(AWSEndpointEnv), "the Endpoint.")

}

func (sc *SessionConfig) Revise() {
	logs.Logger.Infof("AWS SessionConfig is Region: %v; Endpoint: %v", sc.Region, sc.Endpoint)
}
