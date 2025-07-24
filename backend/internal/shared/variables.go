package shared

var (
	/*// storageconfig.go //*/
		ActiveStorage Storage
		ConfigDirPath string
	// storageconfig.go //

	/*// serverpassconfig.go //*/
		ServerPassConfig ServerPassConfigModel
	// serverpassconfig.go //

	/*// metadata.go //*/
		RootDirId string
		FileMetadataMap = make(map[string]*FileMetadata) // file id => FileMetadata
		FolderMetadataMap = make(map[string]*FolderMetadata) // folder id => FolderMetadata
		UserTagsMetadata = make(map[string]*UserTagMetadata) // user tag id => UserTagMetadata
		// holds file ids for each system tag
		SystemTags = map[string]map[string]struct{}{"images": {}, "videos": {}, "audios": {}, "documents": {}} // e.g. image: {id1: {}, id2: {}, ...}
		// ram cache for custom tags list
		UserTagsItems = make(map[string]map[string]string) // user tag id => map[file/folder id => "file"/"folder"]
	// metadata.go //
)