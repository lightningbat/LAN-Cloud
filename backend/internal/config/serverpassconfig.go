package config

import (
	"bufio"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"lan-cloud/internal/shared"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"golang.org/x/crypto/scrypt"
)

var (
	authDirPath              string // absolute path to auth dir
	serverPassConfigFilePath string // absolute path to serverpassconfig file
)

func LoadServerPassConfig(reset string) error {
	if err := setServerPassConfigFilePath(); err != nil {
		return err
	}
	data, err := os.ReadFile(serverPassConfigFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			if err = setServerPassConfigParams(); err != nil { return err }
			return setPassword("ServerPassConfig file does not exist. New Password will be set.")
		}
		return fmt.Errorf("failed to read serverpassconfig file: %v", err)
	}
	err = json.Unmarshal(data, &shared.ServerPassConfig)
	if err != nil {
		if configResetErr := setServerPassConfigParams(); configResetErr != nil { return configResetErr }
		return setPassword(fmt.Sprintf("ServerPassConfig file is invalid. Password will be set. %v", err))
	}

	/* reset password */
	switch reset {
		case "all":
			if err := setServerPassConfigParams(); err != nil { return err } // reset hashing parameters
			return setPassword("Password reset is evoked. New Password will be set.")
		case "pass":
			if err := validateServerPassConfig(false); err != nil { return err }
			return setPassword("Password reset is evoked with old config. New Password will be set.")
	}
	
	// process reached if reset was not evoked and no error occurred while parsing the file
	if validationErr := validateServerPassConfig(true); validationErr != nil {
		if err = setServerPassConfigParams(); err != nil { return err }
		return setPassword(fmt.Sprintf("ServerPassConfig file is invalid. Password will be set. %v", validationErr))
	}
	return nil
}

// validate serverpassconfig file data
func validateServerPassConfig(checkHashes bool) error {
	if _, err := base64.StdEncoding.DecodeString(shared.ServerPassConfig.Salt); err != nil {
		return fmt.Errorf("invalid salt: %v", err)
	}

	if shared.ServerPassConfig.NonceExpiry == 0 {
		return fmt.Errorf("nonce expiry is empty")
	}

	if checkHashes {
		// Scrypt
		if shared.ServerPassConfig.Scrypt.Hash == "" {
			return fmt.Errorf("scrypt hash is empty")
		}
		if _, err := base64.StdEncoding.DecodeString(shared.ServerPassConfig.Scrypt.Hash); err != nil {
			return fmt.Errorf("invalid scrypt hash: %v", err)
		}
	}
	// Scrypt parameters
	if shared.ServerPassConfig.Scrypt.Params.HashLen == 0 {
		return fmt.Errorf("scrypt hash len is empty")
	}
	if shared.ServerPassConfig.Scrypt.Params.Iterations == 0 {
		return fmt.Errorf("scrypt iterations is empty")
	}
	if shared.ServerPassConfig.Scrypt.Params.BlockSize == 0 {
		return fmt.Errorf("scrypt block size is empty")
	}
	if shared.ServerPassConfig.Scrypt.Params.Parallelism == 0 {
		return fmt.Errorf("scrypt parallelism is empty")
	}
	return nil
}

func setPassword(msg string) error {
	if msg != "" { fmt.Println(msg) }
	pass, err := promptPassword()
	if err != nil {
		return err
	}
	if err := setPassHash(pass); err != nil { return err }
	return saveServerPassConfig()
}

func setPassHash(password string ) error {
	// hashing password
	scryptHash, err := scrypt.Key(
		[]byte(password), 
		[]byte(shared.ServerPassConfig.Salt), 
		shared.ServerPassConfig.Scrypt.Params.Iterations, 
		shared.ServerPassConfig.Scrypt.Params.BlockSize, 
		shared.ServerPassConfig.Scrypt.Params.Parallelism, 
		shared.ServerPassConfig.Scrypt.Params.HashLen,
	)
	if err != nil { return err }
	// convert to base64 string
	shared.ServerPassConfig.Scrypt.Hash = base64.StdEncoding.EncodeToString(scryptHash)

	return nil
}

func setServerPassConfigParams() error {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return fmt.Errorf("failed to generate salt: %v", err)
	}
	shared.ServerPassConfig.Salt = base64.StdEncoding.EncodeToString(salt)
	shared.ServerPassConfig.NonceExpiry = 5000 // 5 seconds
	// Scrypt parameters
	shared.ServerPassConfig.Scrypt.Params.HashLen = 32
	shared.ServerPassConfig.Scrypt.Params.Iterations = 16384
	shared.ServerPassConfig.Scrypt.Params.BlockSize = 8
	shared.ServerPassConfig.Scrypt.Params.Parallelism = 1

	return nil
}

func promptPassword() (string, error) {
	for {
		var password string
		fmt.Print("Enter password: ")
		// read user input till enter is pressed
		reader := bufio.NewReader(os.Stdin)
		password, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		password = strings.TrimSpace(password) // trim leading and trailing spaces
		if password != "" {
			// confirm password
			fmt.Print("You entered: '" + password + "' Confirm (y/n):")
			var confirm string
			fmt.Scanln(&confirm)
			if confirm == "y" || confirm == "Y" {
				return password, nil
			}
		}
	}
}

func setAuthDirPath() error {
	// using program's config directory path already loaded by storageconfig
	authDirPath = filepath.Join(shared.ConfigDirPath, "auth")
	err := os.MkdirAll(authDirPath, 0600)
	if err != nil {
		return fmt.Errorf("failed to create auth directory: %v", err)
	}
	// make auth directory hidden on windows
	return hidePath(authDirPath)
}

func setServerPassConfigFilePath() error {
	if err := setAuthDirPath(); err != nil {
		return err
	}
	serverPassConfigFilePath = filepath.Join(authDirPath, "serverpassconfig.json")
	return nil
}

func saveServerPassConfig() error {
	data, err := json.MarshalIndent(shared.ServerPassConfig, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal serverpassconfig: %v", err)
	}
	err = os.WriteFile(serverPassConfigFilePath, data, 0600)
	if err != nil {
		return fmt.Errorf("failed to write serverpassconfig file: %v", err)
	}
	return nil
}

// make path hidden on windows
func hidePath(path string) error {
	if runtime.GOOS == "windows" {
		if err := exec.Command("attrib", "+h", path).Run(); err != nil {
			return fmt.Errorf("failed to set path as hidden: %v", err)
		}
	}
	return nil
}