package media

import (
	"os"
	"strings"

	"../misc/log"

	"../config"
	"./storage/cloudflare"
	"./storage/s3"
	"./storage/s3glacier"
	"./storage/vimeo"
)

func saveImage(fileObject *FileModel) {
	fileName, bucket := s3.UploadToS3(fileObject.BucketMeta["path"])
	fileObject.BucketMeta = map[string]string{
		"key":    fileName,
		"bucket": bucket,
	}
	fileObject.Bucket = BUCKET_AWS_S3
}

func saveVideo(fileObject *FileModel) {
	uploadVimeo := false

	var localPath = fileObject.BucketMeta["path"]
	if uploadVimeo {
		videoID, videoLink := vimeo.Upload(localPath)
		fileObject.BucketMeta = map[string]string{
			"link": videoLink,
			"uri":  videoID,
		}
		fileObject.Bucket = BUCKET_VIMEO
	} else {
		videoID := cloudflare.Upload(localPath)
		fileObject.BucketMeta = map[string]string{
			"uid": videoID,
		}
		fileObject.Bucket = BUCKET_CLOUDFLARE
	}
}

// MoveMediaSafe Moves media to somewhere safe
func (fileObject *FileModel) MoveMediaSafe() {

	var lowLatencyAccessUploaded bool = false

	var localPath = fileObject.BucketMeta["path"]
	if strings.HasPrefix(fileObject.FileType, "image") {
		saveImage(fileObject)
		lowLatencyAccessUploaded = true
	} else if strings.HasPrefix(fileObject.FileType, "video") {
		saveVideo(fileObject)
		lowLatencyAccessUploaded = true
	}

	fileObject.GlacierArchiveID = s3glacier.Upload(localPath)

	if lowLatencyAccessUploaded {
		os.Remove(config.LOCAL_FOLDER + localPath)
	} else {
		log.Fatal("Low latency Upload failed for file ", fileObject)
	}

}
