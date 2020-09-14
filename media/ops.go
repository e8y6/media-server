package media

import (
	"fmt"
	"os"
	"strings"

	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
)

func (fileObject *FileModel) IsImage() bool {
	return (strings.HasPrefix(fileObject.FileType, "image"))
}

func (fileObject *FileModel) ProcessMedia() {

	if fileObject.IsImage() {

		targetWidth := 720

		img, err := imgio.Open("./persist/" + fileObject.Path)
		if err != nil {
			fmt.Println(err)
			return
		}

		srcWidth, srcHeight := img.Bounds().Dx(), img.Bounds().Dy()

		if srcWidth < targetWidth {
			targetWidth = srcWidth
		}

		savePath := (strings.Split(fileObject.Path, "."))[0] + "_opt.jpg"

		targetHeight := int(float64(srcHeight) * (float64(targetWidth) / float64(srcWidth)))
		result := transform.Resize(img, targetWidth, targetHeight, transform.NearestNeighbor)
		if err := imgio.Save("./persist/"+savePath, result, imgio.JPEGEncoder(85)); err != nil {
			fmt.Println(err)
			return
		}
		os.Remove("./persist/" + fileObject.Path)

		fileObject.Path = savePath

		fileObject.Save()

	}

}
