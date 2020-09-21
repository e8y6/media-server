package media

import (
	"fmt"
	"os"
	"strings"

	"../config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	panic(msg)
}

func UploadToS3(localPath string) (string, string) {

	bucket := config.AWSBuckets["media_store"]

	file, err := os.Open("persist/" + localPath)
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

	os.Remove("persist/" + localPath)

	return localPath, bucket

}

func (fileObject *FileModel) MoveMediaSafe() {

	if strings.HasPrefix(fileObject.FileType, "image") {
		fileName, bucket := UploadToS3(fileObject.BucketMeta["path"])
		fileObject.BucketMeta = map[string]string{
			"key":    fileName,
			"bucket": bucket,
		}
		fileObject.Bucket = BUCKET_AWS_S3
	} else {

		fmt.Println("not image...")
	}

}
