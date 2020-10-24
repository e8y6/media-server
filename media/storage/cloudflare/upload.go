package cloudflare

import (
	"net/http"
	"os"
	"strings"

	"../../../config"
	"../../../misc/exceptions"
	"github.com/eventials/go-tus"
)

func Upload(localPath string) string {

	f, err := os.Open(config.LOCAL_FOLDER + localPath)

	if err != nil {
		panic(exceptions.Exception{
			Cause: "Unable to open specified file.",
			Type:  exceptions.TYPE_INTERNAL_ERROR,
			Error: err,
		})
	}

	defer f.Close()

	headers := make(http.Header)
	headers.Add("X-Auth-Email", config.CF_EMAIL)
	headers.Add("X-Auth-Key", config.CF_STREAM_TOKEN)

	mConfig := &tus.Config{
		ChunkSize:           5 * 1024 * 1024, // Cloudflare Stream requires a minimum chunk size of 5MB.
		Resume:              false,
		OverridePatchMethod: false,
		Store:               nil,
		Header:              headers,
		HttpClient:          nil,
	}

	client, err := tus.NewClient("https://api.cloudflare.com/client/v4/accounts/"+config.CF_ACCOUNT_ID+"/stream", mConfig)
	if err != nil {
		panic(exceptions.Exception{
			Cause: "CF: Unable to create new tus client.",
			Type:  exceptions.TYPE_INTERNAL_ERROR,
			Error: err,
		})
	}

	upload, err := tus.NewUploadFromFile(f)
	if err != nil {
		panic(exceptions.Exception{
			Cause: "CF: Unable to create new tus upload from client.",
			Type:  exceptions.TYPE_INTERNAL_ERROR,
			Error: err,
		})
	}

	upload.Metadata["requiresignedurls"] = "true"

	uploader, err := client.CreateUpload(upload)
	if err != nil {
		panic(exceptions.Exception{
			Cause: "CF: Unable to create new tus uploader.",
			Type:  exceptions.TYPE_INTERNAL_ERROR,
			Error: err,
		})
	}

	uploadUrlExploded := strings.Split(uploader.Url(), "/")

	if len(uploadUrlExploded) < 3 {
		// TODO also check for uid regex
		panic(exceptions.Exception{
			Cause: "Unable to fetch UID info for " + localPath,
			Type:  exceptions.TYPE_INTERNAL_ERROR,
			Error: err,
		})
	}

	err = uploader.Upload()
	if err != nil {
		panic(exceptions.Exception{
			Cause: "CF: Upload failed for ." + localPath,
			Type:  exceptions.TYPE_INTERNAL_ERROR,
			Error: err,
		})
	}

	return uploadUrlExploded[len(uploadUrlExploded)-1]
}
