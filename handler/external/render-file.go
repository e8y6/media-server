package external

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"../../config"
	"../../media"
	"../../media/storage/cloudflare"
	"../../media/storage/vimeo"
	"../../misc/exceptions"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gorilla/mux"
)

func serveFromLocal(w *http.ResponseWriter, path string) {
	// TODO copy as stream
	data, err := ioutil.ReadFile(config.LOCAL_FOLDER + path)
	if err != nil {
		panic(exceptions.Exception{
			Type:  exceptions.TYPE_INTERNAL_ERROR,
			Error: err,
		})
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
		panic(exceptions.Exception{
			Type:  exceptions.TYPE_INTERNAL_ERROR,
			Error: err,
			Cause: "Error ocurred while find object in S3",
		})
	}

	_, err = io.Copy(*w, result.Body)
	if err != nil {
		panic(exceptions.Exception{
			Message: "Socket Hangup",
			Type:    exceptions.TYPE_BAD_REQUEST,
			Error:   err,
			Cause:   "Error ocurred copying data from S3 to client fd",
		})
	}

}

func serveFromVimeo(w *http.ResponseWriter, bucketMeta map[string]string) {

	url := vimeo.GetStreamingURL(bucketMeta["uri"])
	op, _ := json.Marshal(map[string]string{
		"url": url,
	})

	(*w).Header().Set("Content-Type", "application/json")
	(*w).Write(op)

}

func serveFromCloudflare(w *http.ResponseWriter, bucketMeta map[string]string) {

	url := cloudflare.GetSignedURL(bucketMeta["uid"])

	op, _ := json.Marshal(map[string]string{
		"url": url,
	})

	(*w).Header().Set("Content-Type", "application/json")
	(*w).Write(op)

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
	} else if result.Bucket == media.BUCKET_CLOUDFLARE {
		serveFromCloudflare(&w, result.BucketMeta)
	}

}
