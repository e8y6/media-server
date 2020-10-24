package s3glacier

import (
	"bytes"
	"io/ioutil"

	"../../../config"
	"../../../misc/exceptions"
	"../../../misc/log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/glacier"
)

// Copied from Somewhere

func Upload(localPath string) string {

	sess, err := session.NewSession(&aws.Config{
		Credentials: config.AWSCredentials,
		Region:      &config.AWSRegion,
	})

	// Create Glacier client in default region
	svc := glacier.New(sess)

	// start snippet
	vaultName := "backup"

	buf, err := ioutil.ReadFile(config.LOCAL_FOLDER + localPath)
	if err != nil {
		panic(exceptions.Exception{
			Cause: "GL: Error while reading local file : " + localPath,
			Type:  exceptions.TYPE_INTERNAL_ERROR,
			Error: err,
		})
	}

	result, err := svc.UploadArchive(&glacier.UploadArchiveInput{
		VaultName: &vaultName,
		Body:      bytes.NewReader(buf),
	})
	if err != nil {
		panic(exceptions.Exception{
			Cause: "GL: Error while uploading file : " + localPath,
			Type:  exceptions.TYPE_INTERNAL_ERROR,
			Error: err,
		})
	}

	log.Info("Uploaded to archive", *result.ArchiveId)

	return *result.ArchiveId
}
