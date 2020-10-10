package vimeo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"../../../config"
)

func UploadToVimeo(localPath string) (videoID string, videoLink string) {

	videoURI, videoLink, uploadURL := createVideo(localPath)
	videoID = strings.Split(videoURI, "/")[2]

	uploadOriginalFile(uploadURL, localPath)
	moveVideo(videoURI, config.VIMEO_Folders["generic"])

	return videoID, videoLink

}

func createVideo(name string) (videoURI string, videoLink string, uploadURL string) {

	url := VIMEO_BURL + "/me/videos"
	method := "POST"

	params := map[string]string{
		"name":         name,
		"privacy.view": "nobody",
	}
	payloadParams, _ := json.Marshal(params)
	payload := strings.NewReader(string(payloadParams))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		panic(err)
	}
	req.Header.Add("Authorization", "Bearer "+config.VIMEO_OAuthToken)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var responseDecodedMap map[string]json.RawMessage
	err = json.Unmarshal(body, &responseDecodedMap)
	if err != nil {
		panic(err)
	}

	var uploadMap map[string]string
	err = json.Unmarshal(responseDecodedMap["upload"], &uploadMap)
	if err != nil {
		panic(err)
	}

	uploadURL = uploadMap["upload_link"]
	json.Unmarshal(responseDecodedMap["uri"], &videoURI)
	json.Unmarshal(responseDecodedMap["link"], &videoLink)

	return videoURI, videoLink, uploadURL

}

func moveVideo(videoURI string, folder string) {

	// url := VIMEO_BURL + "/me/projects/2604204/videos/461715337"
	url := VIMEO_BURL + "/me/projects/" + folder + videoURI
	method := "PUT"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		panic(err)
	}
	req.Header.Add("Authorization", "Bearer "+config.VIMEO_OAuthToken)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer res.Body.Close()
	_, err = ioutil.ReadAll(res.Body)

	if err != nil {
		fmt.Println(err)
	}

	// TODO Handle error here
}

func uploadOriginalFile(uploadURL string, localPath string) {

	url := uploadURL
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	file, errFile1 := os.Open(config.LOCAL_FOLDER + localPath)
	defer file.Close()
	part1,
		errFile1 := writer.CreateFormFile("file_data", filepath.Base(config.LOCAL_FOLDER+localPath))
	_, errFile1 = io.Copy(part1, file)
	if errFile1 != nil {
		panic(errFile1)
	}
	err := writer.Close()
	if err != nil {
		panic(err)
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println("Vimeo: Upload success for ", localPath)
}
