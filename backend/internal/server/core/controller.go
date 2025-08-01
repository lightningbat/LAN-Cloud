package core

import (
	"encoding/json"
	"fmt"
	"lan-cloud/internal/metadata"
	"lan-cloud/internal/shared"
	"lan-cloud/internal/utils"
	"os"
	"path/filepath"
)

func GetFolder(folderId string) (data []byte, err error) {
	folderData, ok := shared.FolderMetadataMap[folderId]
	if !ok {
		return nil, fmt.Errorf("folder %s not found", folderId)
	}
	// get all files in folder
	filesMetadata := make(map[string]*shared.FileMetadata) // fileId => fileName
	for fileId := range folderData.Files {
		fileMetadata, ok := shared.FileMetadataMap[fileId]
		if ok {
			filesMetadata[fileId] = fileMetadata
		}
	}
	// get all subfolders in folder
	subfoldersMetadata := make(map[string]*shared.FolderMetadata) // folderId => folderName
	for folderId := range folderData.SubFolders {
		folderMetadata, ok := shared.FolderMetadataMap[folderId]
		if ok {
			subfoldersMetadata[folderId] = folderMetadata
		}
	}

	resultData := struct {
		Folder     *shared.FolderMetadata `json:"folder"`
		Files      map[string]*shared.FileMetadata `json:"files"`
		Subfolders map[string]*shared.FolderMetadata `json:"subfolders"`
	}{
		Folder:     folderData,
		Files:      filesMetadata,
		Subfolders: subfoldersMetadata,
	}

	responseData, err := json.Marshal(resultData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal folder data: %v", err)
	}
	return responseData, nil
}

func GetTags() (data []byte, err error) {
	type UserTag struct {
		Name  string `json:"name"`
		Color string `json:"color"`
		ItemCount int `json:"item_count"`
	}
	var response = struct {
		ImageCount    int `json:"image_count"`
		VideoCount    int `json:"video_count"`
		AudioCount    int `json:"audio_count"`
		DocumentCount int `json:"document_count"`
		CustomTags    map[string]UserTag `json:"custom_tags"`
	}{CustomTags: make(map[string]UserTag)}
	response.ImageCount = len(shared.SystemTags["images"])
	response.VideoCount = len(shared.SystemTags["videos"])
	response.AudioCount = len(shared.SystemTags["audios"])
	response.DocumentCount = len(shared.SystemTags["documents"])
	for id, tag := range shared.UserTagsMetadata {
		response.CustomTags[id] = UserTag{
			Name:  tag.Name,
			Color: tag.Color,
			ItemCount: len(shared.UserTagsItems[id]),
		}
	}
	data, err = json.Marshal(response)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal tags data: %v", err)
	}
	return data, nil
}

func GetTagItems(tag_type string, tag_id string) (data []byte, err error) {
	if tag_type == "System" {
		return json.Marshal(shared.SystemTags[tag_id])
	} else {
		return json.Marshal(shared.UserTagsItems[tag_id])
	}
}

func GetFilesMetadata(fileIds *[]string) ([]byte, error) {
	filesMetadata := make(map[string]*shared.FileMetadata) // fileId => fileName
	for _, fileId := range *fileIds {
		fileMetadata, ok := shared.FileMetadataMap[fileId]
		if ok {
			filesMetadata[fileId] = fileMetadata
		}
	}
	return json.Marshal(filesMetadata)
}

func GetFoldersMetadata(folderIds *[]string) ([]byte, error) {
	foldersMetadata := make(map[string]*shared.FolderMetadata) // folderId => folderName
	for _, folderId := range *folderIds {
		folderMetadata, ok := shared.FolderMetadataMap[folderId]
		if ok {
			foldersMetadata[folderId] = folderMetadata
		}
	}
	return json.Marshal(foldersMetadata)
}

func Rename(id string, itemType string, newName string) error {
	var old_name string
	if itemType == "folder" {
		old_name = shared.FolderMetadataMap[id].Name
	} else {
		old_name = shared.FileMetadataMap[id].Name
	}

	path := utils.GetAbsolutePath(id, itemType)

	// rename item in file system
	err := os.Rename(filepath.Join(path, old_name), filepath.Join(path, newName))
	if err != nil {
		return fmt.Errorf("failed to rename item: %v", err)
	}

	// update item metadata
	if itemType == "folder" {
		shared.FolderMetadataMap[id].Name = newName
		metadata.SaveFolderMetadata()
	} else {
		shared.FileMetadataMap[id].Name = newName
		metadata.SaveFileMetadata()
	}

	return nil
}