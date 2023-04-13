package awsapi

import (
	"os"

	"github.com/aws/aws-sdk-go/service/storagegateway"
	"github.com/spf13/pflag"

	"zq-xu/helper/config"
	"zq-xu/helper/log"
	"zq-xu/helper/utils"
)

const (
	awsS3ConfigName = "AWSS3Config"

	AWSS3ACLEnv    = "AWSS3ACL"
	AWSS3BucketEnv = "AWSS3Bucket"
	AWSS3VolumeEnv = "AWSS3Volume"
)

var (
	S3Cfg = &S3Config{}
)

type S3Config struct {
	// public read for storagegateway.ObjectACLPublicRead
	ACL string

	Bucket string

	Volume string
}

func InitS3Config() {
	config.RegisterCfg(awsS3ConfigName, S3Cfg)
}

// AddFlags adds flags for S3Config
func (sc *S3Config) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&sc.ACL, "aws-s3-acl", utils.GetStringEnvWithDefault(AWSS3ACLEnv, storagegateway.ObjectACLPublicRead), "the ACL.")
	fs.StringVar(&sc.Bucket, "aws-s3-bucket", os.Getenv(AWSS3BucketEnv), "the bucket.")
	fs.StringVar(&sc.Volume, "aws-s3-volume", os.Getenv(AWSS3VolumeEnv), "the volume.")
}

func (sc *S3Config) Revise() {
	log.Logger.Infof("AWS S3Config is ACL: %v", sc.ACL)
}
