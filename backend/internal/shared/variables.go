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
		FileMetadataMap = make(map[string]*FileMetadata) // file id => FileMetadata
		FolderMetadataMap = make(map[string]*FolderMetadata) // folder id => FolderMetadata
		UserTagsMetadata = make(map[string]*UserTagMetadata) // user tag id => UserTagMetadata
		// for faster lookup for file/folder id by relative path
		FolderRelPathToId = make(map[string]string) // relative path => folder id
		FileRelPathToId  = make(map[string]string) // relative path => file id
		// holds file ids for each system tag
		SystemTags = map[string]map[string]struct{}{"image": {}, "video": {}, "audio": {}, "document": {}} // e.g. image: {id1: {}, id2: {}, ...}
		// ram cache for custom tags list
		UserTagsItems = make(map[string]map[string]string) // user tag id => map[file/folder id => "file"/"folder"]
	// metadata.go //
)