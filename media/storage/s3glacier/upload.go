package s3glacier

import (
	"bytes"
	"io/ioutil"
	"log"

	"../../../config"

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

	result, err := svc.UploadArchive(&glacier.UploadArchiveInput{
		VaultName: &vaultName,
		Body:      bytes.NewReader(buf),
	})
	if err != nil {
		panic(err)
	}

	log.Println("Uploaded to archive", *result.ArchiveId)

	return *result.ArchiveId
}
