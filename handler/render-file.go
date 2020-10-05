package handler

import (
	"io"
	"io/ioutil"
	"net/http"

	"../config"
	"../media"
	"../media/storage/vimeo"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gorilla/mux"
)

func serveFromLocal(w *http.ResponseWriter, path string) {
	// TODO copy as stream
	data, err := ioutil.ReadFile("persist/" + path)
	if err != nil {
		panic(err)
	}
	(*w).Write(data)
}

func serveFromAWSS3(w *http.ResponseWriter, bucketMeta map[string]string) {

	svc := s3.New(session.New(), &aws.Config{
		Credentials: config.AWSCredentials,
		Region:      &config.AWSRegion,
	})

	result, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucketMeta["bucket"]),
		Key:    aws.String(bucketMeta["key"]),
	})
	if err != nil {
		panic("AWS S3 error")
	}

	_, err = io.Copy(*w, result.Body)
	if err != nil {
		panic(err)
	}

}

func serveFromVimeo(w *http.ResponseWriter, bucketMeta map[string]string) {

	data := vimeo.GetVideoFilesAsBytes(bucketMeta["uri"])

	(*w).Header().Set("Content-Type", "application/json")
	(*w).Write(data)

}

// RenderFile renders file with file ID
func RenderFile(w http.ResponseWriter, r *http.Request) {

	urlParams := mux.Vars(r)
	fileID := urlParams["id"]
	result := media.GetFileDetails(fileID)

	if result.Bucket == media.BUCKET_LOCAL {
		serveFromLocal(&w, result.BucketMeta["path"])
	} else if result.Bucket == media.BUCKET_AWS_S3 {
		serveFromAWSS3(&w, result.BucketMeta)
	} else if result.Bucket == media.BUCKET_VIMEO {
		serveFromVimeo(&w, result.BucketMeta)
	}

}
