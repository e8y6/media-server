package media

import (
	"strings"

	"./storage/s3"
	"./storage/vimeo"
)

// MoveMediaSafe Moves media to somewhere safe
func (fileObject *FileModel) MoveMediaSafe() {

	if strings.HasPrefix(fileObject.FileType, "image") {

		fileName, bucket := s3.UploadToS3(fileObject.BucketMeta["path"])
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
