package media

import (
	"fmt"
	"os"
	"strings"

	"../config"
	"golang.org/x/oauth2"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	tus "github.com/eventials/go-tus"
	"github.com/silentsokolov/go-vimeo/vimeo"
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

// Upload to vimeo
type Uploader struct{}

func (u Uploader) UploadFromFile(c *vimeo.Client, uploadURL string, f *os.File) error {
	tusClient, err := tus.NewClient(uploadURL, nil)
	if err != nil {
		return err
	}

	upload, err := tus.NewUploadFromFile(f)
	if err != nil {
		return err
	}

	uploader := tus.NewUploader(tusClient, uploadURL, upload, 0)

	return uploader.Upload()
}

func UploadToVimeo(localPath string) (*vimeo.Video, *vimeo.Response) {

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.VIMEO_OAuthToken},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	config := vimeo.Config{
		Uploader: &Uploader{},
	}

	client := vimeo.NewClient(tc, &config)

	f, _ := os.Open("./persist/" + localPath)

	video, resp, err := client.Users.UploadVideo("", f)

	if err != nil {
		panic(err)
	}

	os.Remove("persist/" + localPath)

	return video, resp

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

		video, _ := UploadToVimeo(fileObject.BucketMeta["path"])
		fileObject.BucketMeta = map[string]string{
			"resource_key": video.ResourceKey,
			"link":         video.Link,
			"uri":          video.URI,
		}
		fileObject.Bucket = BUCKET_VIMEO

	}

}
