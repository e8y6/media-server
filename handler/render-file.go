package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"../media"
	"github.com/gorilla/mux"
)

func RenderFile(w http.ResponseWriter, r *http.Request) {

	urlParams := mux.Vars(r)
	fileID := urlParams["id"]
	result := media.GetFileDetails(fileID)

	mjson, _ := json.Marshal(result)
	fmt.Println(string(mjson))

	if result.Bucket == media.BUCKET_LOCAL {
		fmt.Println("persist/" + result.BucketMeta["path"])
		data, err := ioutil.ReadFile("persist/" + result.BucketMeta["path"])
		if err != nil {
			panic(err)
		}
		w.Write(data)
	}

}
