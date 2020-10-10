package s3

import (
	"os"

	"../../../config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func UploadToS3(localPath string) (string, string) {

	bucket := config.AWSBuckets["media_store"]

	file, err := os.Open(config.LOCAL_FOLDER + localPath)
	if err != nil {
		panic("Unable to open file")
	}
	defer file.Close()

	sess, err := session.NewSession(&aws.Config{
		Credentials: config.AWSCredentials,
		Region:      &config.AWSRegion,
	})

	uploader := s3manager.NewUploader(sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(localPath),
		Body:   file,
	})
	if err != nil {
		panic(err)
	}

	os.Remove(config.LOCAL_FOLDER + localPath)

	return localPath, bucket

}
