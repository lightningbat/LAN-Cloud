package config

import (
	"bufio"
	"encoding/json"
	"fmt"
	"lan-cloud/internal/shared"
	"os"
	"path/filepath"
	"strings"
)

type StorageConfig struct {
	Storages           []shared.Storage `json:"storages"`
	ActiveStorageIndex int       		`json:"active_storage"`
}

var (
	storageConfigFilePath string // absolute path to storageconfig file
	storageConfig  StorageConfig
)

func LoadStorageConfig(cliStoragePath string) error {
	if err := loadConfigFile(); err != nil { return err }
	// if storage path is provided via flag
	cliStoragePath = strings.TrimSpace(cliStoragePath)
	if cliStoragePath != "" {
		if err := updateActiveStorage(cliStoragePath); err != nil { return err }
		finalize()
		return nil
	} else {
		if len(storageConfig.Storages) > storageConfig.ActiveStorageIndex { // if storageconfig has active storage
			if isWritable(storageConfig.Storages[storageConfig.ActiveStorageIndex].Path) { // if active storage path is writable
				finalize()
				return nil
			} else {
				fmt.Println("storage path read from storage config file is not writable: " + storageConfig.Storages[storageConfig.ActiveStorageIndex].Path)
				fmt.Println("falling back to prompt for storage path...")
			}
		}
		storagePath, err := promptForStoragePath()
		if err != nil { return err }
		if err = updateActiveStorage(storagePath); err != nil { return err }
		finalize()
		return nil
	}
}

// called at the end of main
func finalize()  {
	// set globally accessible variables to active storage
	shared.ActiveStorage = storageConfig.Storages[storageConfig.ActiveStorageIndex]
	// release local file memory
	storageConfigFilePath = ""
	storageConfig = StorageConfig{}
}

func promptForStoragePath() (string, error) {
	defaultStorage, err := getDefaultStoragePath() // get system default storage path
	if err != nil { return "", err }
	for {
		var storagePath string
		fmt.Println("Enter storage path or press enter to use default storage path (default: " + defaultStorage + "):")
		// read user input till enter is pressed
		reader := bufio.NewReader(os.Stdin)
		storagePath, err = reader.ReadString('\n')
		if err != nil { return "", err }

		storagePath = strings.TrimSpace(storagePath) // trim leading and trailing spaces
		if storagePath == "" {
			return defaultStorage, nil
		} // if user pressed enter, return default storage path
		if isWritable(storagePath) {
			return storagePath, nil
		} else {
			fmt.Println("storage path is not writable: " + storagePath)
		}
	}
}

func updateActiveStorage(path string) error {
	// if provided storage path is writable
	if !isWritable(path) {
		return fmt.Errorf("storage path is not writable: %s", path)
	}
	// check if provided storage path exists
	for index, storage := range storageConfig.Storages {
		if storage.Path == path {
			storageConfig.ActiveStorageIndex = index
			return saveConfig()
		}
	}
	// if provided storage path doesn't exist in storageconfig
	storage := shared.Storage{Path: path, Metadata: shared.HashPath(path)} // { Path: storagePath, Metadata: hashedMetadataFileContainerName }
	storageConfig.Storages = append(storageConfig.Storages, storage)
	storageConfig.ActiveStorageIndex = len(storageConfig.Storages) - 1
	return saveConfig()
}

// loads storageconfig file into memory(storageconfig)
// creates empty storageconfig file if it doesn't exist
func loadConfigFile() error {
	if err := setConfigDirPath(); err != nil {
		return err
	}
	storageConfigFilePath = filepath.Join(shared.ConfigDirPath, "storageconfig.json")
	data, err := os.ReadFile(storageConfigFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			// create empty storageconfig file
			return saveConfig()
		}
		return fmt.Errorf("failed to read storageconfig file: %v", err)
	}
	err = json.Unmarshal(data, &storageConfig)
	if err != nil {
		return fmt.Errorf("failed to read storageconfig file: %v", err)
	}
	return nil
}

// sets storageconfig directory path, creates storageconfig directory if it doesn't exist
func setConfigDirPath() error {
	dir, err := os.UserConfigDir()
	if err != nil {
		return fmt.Errorf("failed to determine storageconfig directory: %v", err)
	}
	shared.ConfigDirPath = filepath.Join(dir, "LAN Cloud") // set global variable
	err = os.MkdirAll(shared.ConfigDirPath, 0755)
	if err != nil {
		return fmt.Errorf("failed to create storageconfig directory: %v", err)
	}
	return nil
}

// returns system default storage path
func getDefaultStoragePath() (string, error) {
	dir, err := os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf("failed to determine cache directory: %v", err)
	}
	defaultPath := filepath.Join(dir, "LAN Cloud", "storage")
	err = os.MkdirAll(defaultPath, 0755)
	if err != nil {
		return "", fmt.Errorf("failed to create storage directory: %v", err)
	}
	return defaultPath, nil
}

// saves storageconfig to file
func saveConfig() error {
	data, err := json.MarshalIndent(storageConfig, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal storageconfig to JSON: %v", err)
	}
	err = os.WriteFile(storageConfigFilePath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write storageconfig to file %s: %v", storageConfigFilePath, err)
	}
	return nil
}

// checks if provided path is writable
func isWritable(path string) bool {
	testFile := filepath.Join(path, ".testwrite")
	err := os.WriteFile(testFile, []byte("test"), 0644)
	if err != nil {
		return false
	}
	os.Remove(testFile) // clean up
	return true
}
