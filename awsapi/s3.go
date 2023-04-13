// refer to https://blog.csdn.net/Edu_enth/article/details/125993181
// aws docs:
// https://docs.aws.amazon.com/zh_cn/sdk-for-go/v1/developer-guide/sdk-utilities.html
// https://docs.aws.amazon.com/sdk-for-go/api/service/s3/#S3.PutObject
// https://docs.aws.amazon.com/sdk-for-go/api/service/s3/#pkg-examples
package awsapi

import (
	"io"
	"os"
	"path"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var S3Client *s3Client

type s3Client struct {
	session  *session.Session
	s3Client *s3.S3
	uploader *s3manager.Uploader
	cfg      *S3Config
}

// invoked after the session config and the s3 config are initialized
func InitS3Client() {
	S3Client = NewS3Client(MustAWSSession(SessionCfg), S3Cfg)
}

func NewS3Client(sess *session.Session, cfg *S3Config) *s3Client {
	if cfg == nil {
		cfg = &S3Config{}
	}

	cli := &s3Client{
		session: sess,
		cfg:     cfg}

	cli.s3Client = s3.New(cli.session)
	cli.uploader = s3manager.NewUploaderWithClient(cli.s3Client, func(u *s3manager.Uploader) {
		u.PartSize = 64 * 1024 * 1024
		u.Concurrency = s3manager.DefaultUploadConcurrency
	})

	return cli
}

func (as3 *s3Client) PutS3Object(filePath, bucket, bucketPath string) (*s3.PutObjectOutput, error) {
	fOpen, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer fOpen.Close()

	return as3.PutS3ObjectByReader(fOpen, bucket, bucketPath)
}

func (as3 *s3Client) PutS3ObjectByReader(fr io.Reader, bucket, bucketPath string) (*s3.PutObjectOutput, error) {
	input := &s3.PutObjectInput{
		Body:   aws.ReadSeekCloser(fr),
		Bucket: aws.String(bucket),
		Key:    aws.String(bucketPath),
		ACL:    &as3.cfg.ACL,
	}

	return as3.s3Client.PutObject(input)
}

func (as3 *s3Client) UploadFile(filePath, bucket, bucketPath string) (*s3manager.UploadOutput, error) {
	fOpen, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer fOpen.Close()

	return as3.UploadFileByReader(fOpen, bucket, bucketPath)
}

func (as3 *s3Client) UploadFileByReader(fr io.Reader, bucket, bucketPath string) (*s3manager.UploadOutput, error) {
	return as3.uploader.Upload(&s3manager.UploadInput{
		Body:   fr,
		Bucket: aws.String(bucket),
		Key:    aws.String(bucketPath),
		ACL:    &as3.cfg.ACL,
	})
}

func GenerateS3BucketPath(ps ...string) string {
	return path.Join(append([]string{S3Cfg.Volume}, ps...)...)
}
