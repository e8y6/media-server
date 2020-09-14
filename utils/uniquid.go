package utils

import (
	"fmt"
	"strings"
	"time"
)

func getCurrentTime() string {
	s := fmt.Sprintf("%d", time.Now().UnixNano())
	return s
}

// GenerateFileName only has precision upto ns.
func GenerateFileName(currnetName string) string {

	extension := strings.Split(currnetName, ".")[1]
	return getCurrentTime() + "." + extension

}
