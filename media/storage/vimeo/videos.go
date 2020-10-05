package vimeo

import (
	"io/ioutil"
	"net/http"

	"../../../config"
)

func GetVideoFilesAsBytes(videoID string) []byte {

	url := VIMEO_BURL + "/videos/" + videoID + "?fields=files"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		panic(err)
	}
	req.Header.Add("Authorization", "Bearer "+config.VIMEO_OAuthToken)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		panic(err)
	}

	return body
}
