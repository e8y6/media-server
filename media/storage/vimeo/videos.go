package vimeo

import (
	"encoding/json"
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

	return []byte(body)
}

type videoQualities struct {
	Quality string `json:"quality"`
	Link    string `json:"link"`
}

func GetStreamingURL(videoID string) (lastUrl string) {

	videoDetails := GetVideoFilesAsBytes(videoID)

	var UnMJson map[string]json.RawMessage
	json.Unmarshal(videoDetails, &UnMJson)

	var availQualities []videoQualities

	json.Unmarshal(UnMJson["files"], &availQualities)

	for _, quality := range availQualities {

		lastUrl = quality.Link
		if quality.Quality == "hls" {
			break
		}

	}

	return lastUrl

}
