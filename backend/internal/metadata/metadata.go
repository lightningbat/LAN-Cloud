package metadata

import (
	"encoding/json"
	"fmt"
	"lan-cloud/internal/shared"
	"lan-cloud/internal/utils"
	"os"
	"path/filepath"
)

// NOTE: File identity is tracked using relative file paths.
// This means if a file is renamed, moved, or replaced outside this application,
// it will be treated as a deletion and remapped as a new file.
//
// True file tracking across renames or replacements would require OS-level hooks
// (e.g., inotify, USN journal) or a full versioned content hashing system,
// which is out of scope for this lightweight LAN storage server.
//
// This trade-off is intentional and simplifies resync logic while remaining
// reliable for real-world usage.
var (
	metadataDirPath string // root dir for all metadata
	metadataFileContPath string // individual storage metadata file container
	fileMetadataPath string // absolute path to FileMetadata
	folderMetadataPath string // absolute path to FolderMetadata
	userTagsMetaDataPath string // absolute path to UserTagsMetadata
)

func Load() error {
	if err := setMetadataFilesPath(); err != nil { return err }
	// read FileMetadataMap file
	if err := loadMetaData(fileMetadataPath); err != nil { return err }
	if err := loadMetaData(folderMetadataPath); err != nil { return err }
	if err := loadMetaData(userTagsMetaDataPath); err != nil { return err }
	loadRelativePathToId()
	return nil
}

// reads file and load into memory
func loadMetaData(filePath string) (error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// create empty metadata files if they don't exist
			switch filePath {
				case fileMetadataPath:
					return SaveFileMetadata()
				case folderMetadataPath:
					return SaveFolderMetadata()
				case userTagsMetaDataPath:
					return SaveUserTagsMetaData()
			}
		}
		return fmt.Errorf("failed to read metadata file: %v", err)
	}
	switch filePath {
		case fileMetadataPath: // load data into FileMetadataMap
			err = json.Unmarshal(data, &shared.FileMetadataMap)
		case folderMetadataPath: // load data into FolderMetadataMap
			err = json.Unmarshal(data, &shared.FolderMetadataMap)
		case userTagsMetaDataPath: // load data into UserTagsMetadata
			err = json.Unmarshal(data, &shared.UserTagsMetadata)
			// create empty user tag id map in user tag items
			for id := range shared.UserTagsMetadata {
				shared.UserTagsItems[id] = make(map[string]string)
			}
	}
	if err != nil {
		return fmt.Errorf("failed to read metadata file: %v", err)
	}
	return nil
}

// sets root metadata directory
func setMetadataDirPath() error {
	configDir := shared.ConfigDirPath
	metadataDirPath = filepath.Join(configDir, "metadata")
	err := os.MkdirAll(metadataDirPath, 0755)
	if err != nil {
		return fmt.Errorf("failed to create metadata directory: %v", err)
	}
	return nil
}

// sets storage based metadata file container
func setMetadataFilesDirPath(dirName string) error {
	if err := setMetadataDirPath(); err != nil { return err }
	metadataFileContPath= filepath.Join(metadataDirPath, dirName)
	err := os.MkdirAll(metadataFileContPath, 0755) // create metadata file container directory
	if err != nil {
		return fmt.Errorf("failed to create metadata file container directory: %v", err)
	}
	return nil
}

// loads storage FileMetadataMap path from application config file
func setMetadataFilesPath() error {
	if err := setMetadataFilesDirPath(shared.ActiveStorage.Metadata); err != nil { return err }
	fileMetadataPath = filepath.Join(metadataFileContPath, "filemetadata.json")
	folderMetadataPath = filepath.Join(metadataFileContPath, "foldermetadata.json")
	userTagsMetaDataPath = filepath.Join(metadataFileContPath, "usertags.json")
	return nil
}

func SaveFileMetadata() error {
	data, err := json.MarshalIndent(shared.FileMetadataMap, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal FileMetadataMap to JSON: %v", err)
	}
	if err := os.WriteFile(fileMetadataPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write FileMetadataMap to file %s: %v", fileMetadataPath, err)
	}
	return nil
}
func SaveFolderMetadata() error {
	data, err := json.MarshalIndent(shared.FolderMetadataMap, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal FolderMetadataMap to JSON: %v", err)
	}
	if err := os.WriteFile(folderMetadataPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write FolderMap to file %s: %v", folderMetadataPath, err)
	}
	return nil
}

func SaveUserTagsMetaData() error {
	data, err := json.MarshalIndent(shared.UserTagsMetadata, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal UserTags to JSON: %v", err)
	}
	if err := os.WriteFile(userTagsMetaDataPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write UserTags to file %s: %v", userTagsMetaDataPath, err)
	}
	return nil
}

func AddFileMetadata(filemetadata *shared.FileMetadata) (id string) {
	id = shared.HashPath(filemetadata.RelativePath) // generate file id from relative path
	// map file metadata id to FilesMetadata list
	shared.FileMetadataMap[id] = filemetadata
	// map relative path to file id
	shared.FileRelPathToId[filemetadata.RelativePath] = id
	// add file id to parent folder's file list
	shared.FolderMetadataMap[filemetadata.ParentId].Files[id] = struct{}{}
	filemetadata.SystemTag = utils.GetFileCategory(filemetadata.Name) // set system tag
	// add file id to it's associate system tag list
	if (filemetadata.SystemTag != "") { shared.SystemTags[filemetadata.SystemTag][id] = struct{}{} }
	return
}
func AddFolderMetadata(foldermetadata *shared.FolderMetadata) (id string) {
	if foldermetadata.RelativePath == "" {
		id = "root" // set root folder id to "root"
	} else {
		id = shared.HashPath(foldermetadata.RelativePath) // generate folder id from relative path
		// add folder id to parent folder subfolder list
		shared.FolderMetadataMap[foldermetadata.ParentId].SubFolders[id] = struct{}{}
	}
	// map folder metadata id to FolderMetadata list
	shared.FolderMetadataMap[id] = foldermetadata
	// map relative path to folder id
	shared.FolderRelPathToId[foldermetadata.RelativePath] = id
	return
}

func DeleteFileMetadata(id string) {
	// delete from relative path to id map
	delete(shared.FileRelPathToId, shared.FileMetadataMap[id].RelativePath)
	// delete id from custom tags list
	for _, tag := range shared.FileMetadataMap[id].Tags {
		delete(shared.UserTagsItems[tag], id)
	}
	// delete id from system tags list
	delete(shared.SystemTags[shared.FileMetadataMap[id].SystemTag], id)
	// delete id from parent folder file list
	delete(shared.FolderMetadataMap[shared.FileMetadataMap[id].ParentId].Files, id)
	// delete from FileMetadataMap
	delete(shared.FileMetadataMap, id)
}

func DeleteFolderMetadata(id string) {
	// delete from relative path to id map
	delete(shared.FolderRelPathToId, shared.FolderMetadataMap[id].RelativePath)
	// delete id from custom tags list
	for _, tag := range shared.FolderMetadataMap[id].Tags {
		delete(shared.UserTagsItems[tag], id)
	}
	// delete from parent folder subfolder list
	delete(shared.FolderMetadataMap[shared.FolderMetadataMap[id].ParentId].SubFolders, id)
	// delete from FolderMetadataMap
	delete(shared.FolderMetadataMap, id)
}

// maps relative path to file/folder id
func loadRelativePathToId() {
	// map relative path to file id
	for id, metaData := range shared.FileMetadataMap {
		shared.FileRelPathToId[metaData.RelativePath] = id
		// add file id to it's associate system tag list
		if (metaData.SystemTag != "") { shared.SystemTags[metaData.SystemTag][id] = struct{}{} }
		// add file id to custom tags list
		for _, tag := range metaData.Tags {
			shared.UserTagsItems[tag][id] = "file"
		}
	}
	// map relative path to folder id
	for id, metaData := range shared.FolderMetadataMap {
		shared.FolderRelPathToId[metaData.RelativePath] = id
		// add folder id to custom tags list
		for _, tag := range metaData.Tags {
			shared.UserTagsItems[tag][id] = "folder"
		}
	}
}