package s3

import (
	"os"

	"../../../config"
	"../../../misc/exceptions"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func UploadToS3(localPath string) (string, string) {

	bucket := config.AWSBuckets["media_store"]

	file, err := os.Open(config.LOCAL_FOLDER + localPath)
	if err != nil {
		panic(exceptions.Exception{
			Cause: "S3: Unable to open local file for upload path : " + localPath,
			Type:  exceptions.TYPE_INTERNAL_ERROR,
			Error: err,
		})
	}
	defer file.Close()

	sess, err := session.NewSession(&aws.Config{
		Credentials: config.AWSCredentials,
		Region:      &config.AWSRegion,
	})

	if err != nil {
		panic(exceptions.Exception{
			Cause: "S3: Error while creating new upload session for file : " + localPath,
			Type:  exceptions.TYPE_INTERNAL_ERROR,
			Error: err,
		})
	}

	uploader := s3manager.NewUploader(sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(localPath),
		Body:   file,
	})
	if err != nil {
		panic(exceptions.Exception{
			Cause: "S3: Error while uploading the file : " + localPath,
			Type:  exceptions.TYPE_INTERNAL_ERROR,
			Error: err,
		})
	}

	return localPath, bucket

}
