package awsapi

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/storagegateway"
	. "github.com/smartystreets/goconvey/convey"
)

// go test -v \
// session.go session_config.go \
// s3_config.go s3.go s3_test.go \
// -test.run TestS3UploadFile -count=1
func TestS3UploadFile(t *testing.T) {
	filePath := "./s3.go"
	bucket := "xzq-bucket"
	bucketPath := "s3.go"

	Convey("TestS3UploadFile", t, func() {
		sessionCfg := &SessionConfig{
			AccessID:     "",
			AccessSecret: "",
			Region:       "eu-central-1",
		}
		s3Cfg := &S3Config{ACL: storagegateway.ObjectACLPublicRead}
		cli := NewS3Client(MustAWSSession(sessionCfg), s3Cfg)

		Convey("PutS3Object", func() {
			_, err := cli.PutS3Object(filePath, bucket, fmt.Sprintf("put-obj-%s", bucketPath))
			So(err, ShouldBeNil)
		})

		Convey("print schema", func() {
			_, err := cli.UploadFile(filePath, bucket, fmt.Sprintf("upload-file-%s", bucketPath))
			So(err, ShouldBeNil)
		})
	})
}
