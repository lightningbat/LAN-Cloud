package filesystem

import (
	"fmt"
	"lan-cloud/internal/metadata"
	"lan-cloud/internal/shared"
	"os"
	"path/filepath"
)

type syncResult struct {
	// list of scanned files
	scannedFolderIds map[string]struct{} // key => folder id
	scannedFileIds map[string]struct{} // key => file id
	newFolderCount int
	newFileCount   int
	fileUpdatedCount int
	folderUpdatedCount int
}

func SyncMetadata() error {
	result := &syncResult{
		scannedFolderIds: make(map[string]struct{}),
		scannedFileIds: make(map[string]struct{}),
	}
	_, _, err := scanDir(shared.ActiveStorage.Path, "root", "", result)
	if err != nil { return err }

	deletedFoldersCount, deletedFilesCount := deleteUntrackedMetadata(&result.scannedFolderIds, &result.scannedFileIds)
	if (result.newFileCount > 0 || deletedFilesCount > 0 || result.fileUpdatedCount > 0){
		if err := metadata.SaveFileMetadata(); err != nil { return err }
	}
	// also include file delete count in folder update cause file ids are also deleted from parent folder metadata
	if (result.newFolderCount > 0 || deletedFoldersCount > 0 || deletedFilesCount > 0 || result.folderUpdatedCount > 0){
		if err := metadata.SaveFolderMetadata(); err != nil { return err }
	}

	fmt.Println("Sync Details:")
	fmt.Printf("\tNew Folders: %d\n", result.newFolderCount)
	fmt.Printf("\tNew Files: %d\n", result.newFileCount)
	fmt.Printf("\tUpdated Files: %d\n", result.fileUpdatedCount)
	fmt.Printf("\tUpdated Folders: %d\n", result.folderUpdatedCount)
	fmt.Printf("\tDeleted Folders: %d\n", deletedFoldersCount)
	fmt.Printf("\tDeleted Files: %d\n", deletedFilesCount)
	return nil
}

func scanDir(absPath string, name string, parentId string, result *syncResult) (childFolderId string, contentSize int64, err error) {

	/** Parent Folder Setup **/
	relativePath := shared.GetRelativePath(absPath)

	folderInfoUpdated := false // indicates if folder metadata was updated

	folderInfo, err := os.Stat(absPath)
	if err != nil { return "", 0, err }
	modifiedTime := folderInfo.ModTime().Unix()

	// Create a folder pointer in the function scope to allow adding files and subfolders 
	// regardless of whether the folder already exists in the metadata.
	var folder *shared.FolderMetadata
	folderId, ok := shared.FolderRelPathToId[relativePath] // check if folder exists in metadata
	if ok {
		folder = shared.FolderMetadataMap[folderId] // get folder metadata from map
		// update folder metadata if modified time is newer
		if (folder.ModifiedTime < modifiedTime && folderId != "root") { // excluding root since a temporary file is created every time to test if root is writable
			folder.ModifiedTime = modifiedTime
			folderInfoUpdated = true
		}
	}else {
		folderId, folder = createFolderMetadata(name, relativePath, parentId, modifiedTime) // create new folder metadata
	}
	result.scannedFolderIds[folderId] = struct{}{}// confirm scanned folder
	
	/** Child Entries Setup **/
	entries, err := os.ReadDir(absPath)
	if err != nil { return "", 0, err }

	for _, entry := range entries {
		entryRelPath := filepath.Join(relativePath, entry.Name())
		if entry.IsDir() {
			// recursive call
			childFolderId, content_size, err := scanDir(filepath.Join(absPath, entry.Name()), entry.Name(), folderId, result)
			if err != nil { return "", 0, err }

			contentSize += content_size // increment child size to parent total size

			if _, ok := folder.SubFolders[childFolderId]; !ok {
				folder.SubFolders[childFolderId] = struct{}{}
				result.newFolderCount++
			}
		} else {
			fileInfo, err := entry.Info()
			if err != nil { return "", 0, err }
			fileSize := fileInfo.Size()
			fileModifiedTime := fileInfo.ModTime().Unix()

			contentSize += fileSize

			fileId, ok := shared.FileRelPathToId[entryRelPath] // check if file metadata exists

			if ok { // if the file exists in metadata
				result.scannedFileIds[fileId] = struct{}{} // confirm scanned file
				// check if file id exists in parent folder's file list
				if _, ok := folder.Files[fileId]; !ok {
					// push file id to parent folder's file list
					folder.Files[fileId] = struct{}{}
					// increment new file count
					result.newFileCount++
				}

				infoUpdated := false // indicates if file metadata was updated
				// sync file size
				if fileSize != shared.FileMetadataMap[fileId].Size {
					shared.FileMetadataMap[fileId].Size = fileSize
					infoUpdated = true
				}
				// sync file modified time
				if fileModifiedTime != shared.FileMetadataMap[fileId].ModifiedTime {
					shared.FileMetadataMap[fileId].ModifiedTime = fileModifiedTime
					infoUpdated = true
				}
				if infoUpdated { result.fileUpdatedCount++ }
			} else { // if the file doesn't exist in metadata
				fileId = createFileMetadata(entry.Name(), entryRelPath, folderId, fileSize, fileModifiedTime) // create new file metadata
				result.scannedFileIds[fileId] = struct{}{} // confirm scanned file
				// push file id to parent folder's file list
				folder.Files[fileId] = struct{}{}
				result.newFileCount++
			}
		}
	}
	if contentSize != folder.Size {
		folder.Size = contentSize
		folderInfoUpdated = true
	}
	if folderInfoUpdated { result.folderUpdatedCount++ }
	return folderId, contentSize, nil
}

func createFolderMetadata(folderName string, folderRelativePath string, parentId string, folderModifiedTime int64) (folderId string, folder *shared.FolderMetadata) {
	folder = &shared.FolderMetadata{
		Name: folderName,
		RelativePath: folderRelativePath,
		ParentId: parentId,
		ModifiedTime: folderModifiedTime,
		Files: make(map[string]struct{}),
		SubFolders: make(map[string]struct{}),
		Owners: []string{},
		Tags: []string{},
	}
	folderId = metadata.AddFolderMetadata(folder)
	return
}

func createFileMetadata(fileName string, fileRelativePath string, parentId string, fileSize int64, fileModifiedTime int64) (fileId string) {
	file := &shared.FileMetadata{
		Name: fileName,
		RelativePath: fileRelativePath,
		ParentId: parentId,
		Size: fileSize,
		ModifiedTime: fileModifiedTime,
		Owners: []string{},
		Tags: []string{},
	}
	fileId = metadata.AddFileMetadata(file)
	return
}

func deleteUntrackedMetadata(scannedFolderIds *map[string]struct{}, scannedFileIds *map[string]struct{}) (deletedFoldersCount int, deletedFilesCount int) {
	deletedFoldersCount = 0
	deletedFilesCount = 0

	// delete metadata of untracked folders
	for folderId, _ := range shared.FolderMetadataMap {
		if _, ok := (*scannedFolderIds)[folderId]; !ok {
			metadata.DeleteFolderMetadata(folderId)
			deletedFoldersCount++
		}
	}
	// delete metadata of untracked files
	for fileId, _ := range shared.FileMetadataMap {
		if _, ok := (*scannedFileIds)[fileId]; !ok {
			metadata.DeleteFileMetadata(fileId)
			deletedFilesCount++
		}
	}
	return
}