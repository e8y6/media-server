package media

import (
	"os"
	"strings"

	"../config"
	"./storage/vimeo"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

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

// Endo of Upload to vimeo

// MoveMediaSafe Moves media to somewhere safe
func (fileObject *FileModel) MoveMediaSafe() {

	if strings.HasPrefix(fileObject.FileType, "image") {

		fileName, bucket := UploadToS3(fileObject.BucketMeta["path"])
		fileObject.BucketMeta = map[string]string{
			"key":    fileName,
			"bucket": bucket,
		}
		fileObject.Bucket = BUCKET_AWS_S3

	} else if strings.HasPrefix(fileObject.FileType, "video") {
		videoID, videoLink := vimeo.UploadToVimeo(fileObject.BucketMeta["path"])
		fileObject.BucketMeta = map[string]string{
			"link": videoLink,
			"uri":  videoID,
		}
		fileObject.Bucket = BUCKET_VIMEO
	}

}
