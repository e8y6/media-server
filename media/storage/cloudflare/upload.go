package cloudflare

import (
	"net/http"
	"os"
	"strings"

	"../../../config"
	"github.com/eventials/go-tus"
)

func Upload(localPath string) string {

	f, err := os.Open(config.LOCAL_FOLDER + localPath)

	if err != nil {
		panic(err)
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

	client, _ := tus.NewClient("https://api.cloudflare.com/client/v4/accounts/"+config.CF_ACCOUNT_ID+"/stream", mConfig)

	upload, _ := tus.NewUploadFromFile(f)

	upload.Metadata["requiresignedurls"] = "true"

	uploader, _ := client.CreateUpload(upload)

	uploadUrlExploded := strings.Split(uploader.Url(), "/")

	uploader.Upload()

	return uploadUrlExploded[len(uploadUrlExploded)-1]
}
