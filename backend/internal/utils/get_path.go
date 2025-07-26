package utils

import (
	"lan-cloud/internal/shared"
	"path/filepath"
)

func GetAbsolutePath(itemId string, itemType string) string {
	var currentParentId string
	switch itemType {
		case "folder":
			currentParentId = shared.FolderMetadataMap[itemId].ParentId
		case "file":
			currentParentId = shared.FileMetadataMap[itemId].ParentId
	}
	segments := recurse(currentParentId, []string{}) // build recursive path
	segments = append([]string{shared.ActiveStorage.Path}, segments...) // add storage path
	return filepath.Join(segments...)
}


func recurse(id string, segments []string) ([]string) {
	if (id == shared.RootDirId) {
		return segments
	}
	segments = append([]string{shared.FolderMetadataMap[id].Name}, segments...) // add folder name
	segments = recurse(shared.FolderMetadataMap[id].ParentId, segments) // recurse
	return segments
}