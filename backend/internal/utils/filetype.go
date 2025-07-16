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

	/* Note: using plural form of file categories, to be consistent with the category names */
	if _, ok := supportedImageExtensions[fileExtension]; ok {
		return "images"
	} else if _, ok := supportedVideoExtensions[fileExtension]; ok {
		return "videos"
	} else if _, ok := supportedAudioExtensions[fileExtension]; ok {
		return "audios"
	} else if _, ok := supportedDocumentExtensions[fileExtension]; ok {
		return "documents"
	} else {
		return ""
	}
}
