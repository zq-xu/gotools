package awsapi

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

func MustAWSSession(cfg *SessionConfig) *session.Session {
	return session.Must(NewAWSSession(cfg))
}

func NewAWSSession(cfg *SessionConfig) (*session.Session, error) {
	if cfg == nil {
		return nil, fmt.Errorf("empty aws session config")
	}

	return session.NewSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials(cfg.AccessID, cfg.AccessSecret, ""),
		Endpoint:         aws.String(cfg.Endpoint),
		Region:           aws.String(cfg.Region),
		S3ForcePathStyle: aws.Bool(false),
	})
}
