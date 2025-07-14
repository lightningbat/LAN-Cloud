package core

import (
	"encoding/json"
	"fmt"
	"lan-cloud/internal/shared"
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
	var response struct {
		ImageCount    int `json:"image_count"`
		VideoCount    int `json:"video_count"`
		AudioCount    int `json:"audio_count"`
		DocumentCount int `json:"document_count"`
		CustomTags map[string]UserTag `json:"custom_tags"`
	}
	response.ImageCount = len(shared.SystemTags["image"])
	response.VideoCount = len(shared.SystemTags["video"])
	response.AudioCount = len(shared.SystemTags["audio"])
	response.DocumentCount = len(shared.SystemTags["document"])
	for id, tag := range shared.UserTagsMetadata {
		response.CustomTags[tag.Name] = UserTag{
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

// func GetFiles(fileIds *[]string) ([]byte, error) {
// 	files := make(map[string]string) // fileId => fileName
// 	for _, fileId := range *fileIds {
// 		file, ok := shared.FileMetadataMap[fileId]
// 		if ok {
// 			files[fileId] = file.Name
// 		}
// 	}
// 	jsonData, err := json.Marshal(files)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to marshal files data: %v", err)
// 	}
// 	return jsonData, nil
// }
