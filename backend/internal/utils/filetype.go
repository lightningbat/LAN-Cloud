package utils

import (
	"strings"
)

var (
	supportedImageExtensions    = map[string]struct{}{"jpg": {}, "jpeg": {}, "png": {}, "gif": {}, "webp": {}, "svg": {}}
	supportedVideoExtensions    = map[string]struct{}{"mp4": {}, "webm": {}}
	supportedAudioExtensions    = map[string]struct{}{"mp3": {}, "m4a": {}, "ogg": {}, "wav": {}}
	supportedDocumentExtensions = map[string]struct{}{"pdf": {}, "txt": {}, "html": {}, "json": {}, "csv": {}, "xm": {}}
)

func GetFileCategory(fileName string) string {
	// check if file name has extension
	if !strings.Contains(fileName, ".") {
		return ""
	}
	fileExtension := strings.Split(fileName, ".")[1]

	if _, ok := supportedImageExtensions[fileExtension]; ok {
		return "image"
	} else if _, ok := supportedVideoExtensions[fileExtension]; ok {
		return "video"
	} else if _, ok := supportedAudioExtensions[fileExtension]; ok {
		return "audio"
	} else if _, ok := supportedDocumentExtensions[fileExtension]; ok {
		return "document"
	} else {
		return ""
	}
}
