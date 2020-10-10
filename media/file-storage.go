package media

import (
	"os"
	"strings"

	"../config"
	"./storage/s3"
	"./storage/s3glacier"
	"./storage/vimeo"
)

// MoveMediaSafe Moves media to somewhere safe
func (fileObject *FileModel) MoveMediaSafe() {

	var lowLatencyAccessUploaded bool = false
	var archiveAccessUploaded bool = false

	var localPath = fileObject.BucketMeta["path"]

	if strings.HasPrefix(fileObject.FileType, "image") {

		fileName, bucket := s3.UploadToS3(localPath)
		fileObject.BucketMeta = map[string]string{
			"key":    fileName,
			"bucket": bucket,
		}
		fileObject.Bucket = BUCKET_AWS_S3
		lowLatencyAccessUploaded = true

	} else if strings.HasPrefix(fileObject.FileType, "video") {
		videoID, videoLink := vimeo.UploadToVimeo(localPath)
		fileObject.BucketMeta = map[string]string{
			"link": videoLink,
			"uri":  videoID,
		}
		fileObject.Bucket = BUCKET_VIMEO
		lowLatencyAccessUploaded = true
	}

	fileObject.GlacierArchiveID = s3glacier.Upload(localPath)
	archiveAccessUploaded = true

	if archiveAccessUploaded && lowLatencyAccessUploaded {
		os.Remove(config.LOCAL_FOLDER + localPath)
	}

}
