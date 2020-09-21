package media

import (
	"fmt"
	"os"

	"../config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	panic(msg)
}

func UploadToS3(localPath string) {

	bucket := "uat-jw-storage"
	filename := localPath

	file, err := os.Open("persist/" + filename)
	if err != nil {
		exitErrorf("Unable to open file %q, %v", err)
	}
	defer file.Close()

	sess, err := session.NewSession(&aws.Config{
		Credentials: config.AWSCredentials,
		Region:      &config.AWSRegion,
	})

	uploader := s3manager.NewUploader(sess)
	x, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(config.AWSBuckets["media_store"]),
		Key:    aws.String(filename),
		Body:   file,
	})
	if err != nil {
		exitErrorf("Unable to upload %q to %q, %v", filename, bucket, err)
	}

	fmt.Println(x)

	fmt.Printf("Successfully uploaded %q to %q\n", filename, bucket)

	fmt.Println("File uploaded...")

}
