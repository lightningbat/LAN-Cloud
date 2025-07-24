package shared

type Storage struct {
	Path     string `json:"path"`
	Metadata string `json:"metadata"`
}

type FileMetadata struct {
	Name         string   `json:"name"`
	ParentId     string   `json:"parent_id"`
	Size         int64    `json:"size"`
	ModifiedTime int64    `json:"modified_time"`
	Owners       []string `json:"owners"`
	Tags         []string `json:"tags"`
	SystemTag    string   `json:"system_tag"` // file type (images, videos, audios, documents, etc.)
}

type FolderMetadata struct {
	Name         string              `json:"name"`
	ParentId     string              `json:"parent_id"`
	Size         int64               `json:"size"`
	ModifiedTime int64               `json:"modified_time"`
	Files        map[string]struct{} `json:"files"`       // key: id
	SubFolders   map[string]struct{} `json:"sub_folders"` // key: id
	Owners       []string            `json:"owners"`
	Tags         []string            `json:"tags"`
}

type UserTagMetadata struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

type ScryptParams struct {
	HashLen     int `json:"hash_len"`
	Iterations  int `json:"iterations"`
	BlockSize   int `json:"block_size"`
	Parallelism int `json:"parallelism"`
}

type ServerPassConfigModel struct {
	Salt        string `json:"salt"`
	NonceExpiry int    `json:"nonce_expiry"`
	Scrypt struct {
		Hash   string       `json:"hash"`
		Params ScryptParams `json:"params"`
	} `json:"scrypt"`
}
