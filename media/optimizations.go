package media

import (
	"fmt"
	"os"
	"strings"

	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
)

func optimizeImage(fileObject *FileModel) {

	targetWidth := 720

	img, err := imgio.Open("./persist/" + fileObject.BucketMeta["path"])
	if err != nil {
		fmt.Println(err)
		return
	}

	srcWidth, srcHeight := img.Bounds().Dx(), img.Bounds().Dy()

	if srcWidth < targetWidth {
		targetWidth = srcWidth
	}

	savePath := (strings.Split(fileObject.BucketMeta["path"], "."))[0] + "_opt.jpeg"

	targetHeight := int(float64(srcHeight) * (float64(targetWidth) / float64(srcWidth)))
	result := transform.Resize(img, targetWidth, targetHeight, transform.NearestNeighbor)
	if err := imgio.Save("./persist/"+savePath, result, imgio.JPEGEncoder(85)); err != nil {
		fmt.Println(err)
		return
	}
	os.Remove("./persist/" + fileObject.BucketMeta["path"])

	fileObject.FileType = "image/jpeg"
	fileObject.BucketMeta["path"] = savePath

	fileObject.Save()

}

// Optimize optimizes the file
func (fileObject *FileModel) Optimize() {

	if strings.HasPrefix(fileObject.FileType, "image") {
		optimizeImage(fileObject)
	}

}