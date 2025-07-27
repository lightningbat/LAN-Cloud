package httpserver

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"lan-cloud/internal/server/core"
	"lan-cloud/internal/server/crypto"
	"lan-cloud/internal/shared"
	"net"
	"net/http"
	"time"
)

/*  /// Authentication */

func handShakeConnect(w http.ResponseWriter, r *http.Request) {
	/* Rate limit */
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	rateLimit_key := ip+"handshake_connect"
	if !getRateLimiter(rateLimit_key, 5, 10*time.Second) { // 5 requests per 10 seconds
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		return
	}

	var requestData struct {
		ClientData string `json:"client_data"`
	}
	err = json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sessionId, nonceStr, err := core.GenerateNonce() // generate nonce
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// convert base64 string to byte array
	hashPassByte, err := base64.StdEncoding.DecodeString(shared.ServerPassConfig.Scrypt.Hash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// aes-gcm encrypt nonce and clientData
	ivBase64, ciphertextBase64, err := crypto.EncryptAESGCM(hashPassByte, []byte(nonceStr+requestData.ClientData))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// write response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Salt                string `json:"salt"`
		SessionId           string `json:"session_id"`
		shared.ScryptParams `json:"scrypt_params"`
		IVBase64            string `json:"iv_base64"`
		CiphertextBase64    string `json:"ciphertext_base64"`
	}{
		shared.ServerPassConfig.Salt,
		sessionId,
		shared.ServerPassConfig.Scrypt.Params,
		ivBase64,
		ciphertextBase64,
	})
}

func handShakeAuthenticate(w http.ResponseWriter, r *http.Request) {
	/* Rate limit */
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	rateLimit_key := ip + "handshake_authenticate"
	if !getRateLimiter(rateLimit_key, 3, 10*time.Second) { // 3 requests per 10 seconds
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		return
	}

	// parse request body
	var requestData struct {
		SessionId string `json:"session_id"`
		AuthKey   string `json:"auth_key"`
	}
	err = json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	password := shared.ServerPassConfig.Scrypt.Hash // use scrypt hash as password
	nonceStr, found := core.NonceCache.Get(requestData.SessionId) // get nonce from session
	if !found {
		http.Error(w, "Invalid session ID", http.StatusBadRequest)
		return
	}

	// convert password and nonce from base64 to byte array
	passwordBytes, err := base64.StdEncoding.DecodeString(password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Use HMAC to hash the password and nonce_str
	mac := hmac.New(sha256.New, passwordBytes)
	mac.Write([]byte(nonceStr.(string)))
	hashedData := mac.Sum(nil)

	// convert hashed data to base64
	hashedDataStr := base64.StdEncoding.EncodeToString(hashedData)

	if hashedDataStr != requestData.AuthKey {
		http.Error(w, "Invalid authentication key", http.StatusUnauthorized)
		return
	}

	// generate session key as new encryption key
	sessionId, sessionKey, err := core.GenerateSessionKey()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// encrypt sessionKey
	ivBase64, ciphertextBase64, err := crypto.EncryptAESGCM(passwordBytes, []byte(sessionKey))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// write response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		SessionId           string `json:"session_id"`
		IVBase64            string `json:"iv_base64"`
		CiphertextBase64    string `json:"ciphertext_base64"`
	}{
		sessionId,
		ivBase64,
		ciphertextBase64,
	})
}

func refreshSessionKey(w http.ResponseWriter, r *http.Request) {
	/* Rate limit */
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	rateLimit_key := ip + "refresh_session_key"
	if !getRateLimiter(rateLimit_key, 10, 60*time.Second) { // 10 requests per minute
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		return
	}

	// parse request body
	var requestData struct {
		SessionId string `json:"session_id"`
	}
	err = json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// get session key from session
	oldSessionKey, found := core.SessionKeyCache.Get(requestData.SessionId)
	if !found {
		http.Error(w, "Invalid session ID", http.StatusBadRequest)
		return
	}

	// generate session key as new encryption key
	sessionId, newSessionKey, err := core.GenerateSessionKey()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// encrypt sessionKey
	ivBase64, ciphertextBase64, err := crypto.EncryptAESGCM(oldSessionKey.([]byte), []byte(newSessionKey))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// write response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		SessionId           string `json:"session_id"`
		IVBase64            string `json:"iv_base64"`
		CiphertextBase64    string `json:"ciphertext_base64"`
	}{
		sessionId,
		ivBase64,
		ciphertextBase64,
	})
}

/* Authentication /// */

func getRootDirId(w http.ResponseWriter, r *http.Request) {
	// rate limit
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	rateLimit_key := ip + "get_root_dir_path"
	if !getRateLimiter(rateLimit_key, 50, 60*time.Second) { // 50 requests per minute
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		RootDirId string `json:"root_dir_id"`
	}{
		RootDirId: shared.RootDirId,
	})
}

func getFolder(w http.ResponseWriter, r *http.Request) {
	// rate limit
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	rateLimit_key := ip + "get_folder"
	if !getRateLimiter(rateLimit_key, 50, 60*time.Second) { // 50 requests per minute
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		return
	}

	// parse request body
	var folder struct {
		FolderId string `json:"folder_id"`
		SessionId string `json:"session_id"`
	} // request body struct
	err = json.NewDecoder(r.Body).Decode(&folder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// get session key from session
	sessionKey, found := core.SessionKeyCache.Get(folder.SessionId)
	if !found {
		http.Error(w, "Invalid session ID", http.StatusBadRequest)
		return
	}

	// get folder data
	data, err := core.GetFolder(folder.FolderId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// encrypt folder data
	ivBase64, ciphertextBase64, err := crypto.EncryptAESGCM(sessionKey.([]byte), data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		IVBase64            string `json:"iv_base64"`
		CiphertextBase64    string `json:"ciphertext_base64"`
	}{
		ivBase64,
		ciphertextBase64,
	})
}

// return basic tags info (tag name, id, color, etc.)
func getTags(w http.ResponseWriter, r *http.Request) {
	// rate limit
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	rateLimit_key := ip + "get_tags"
	if !getRateLimiter(rateLimit_key, 20, 60*time.Second) { // 20 requests per minute
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		return
	}

	// parse request body
	var folder struct {
		SessionId string `json:"session_id"`
	} // request body struct
	err = json.NewDecoder(r.Body).Decode(&folder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// get session key from session
	sessionKey, found := core.SessionKeyCache.Get(folder.SessionId)
	if !found {
		http.Error(w, "Invalid session ID", http.StatusBadRequest)
		return
	}

	// get tags data
	data, err := core.GetTags()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// encrypt tags data
	ivBase64, ciphertextBase64, err := crypto.EncryptAESGCM(sessionKey.([]byte), data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		IVBase64            string `json:"iv_base64"`
		CiphertextBase64    string `json:"ciphertext_base64"`
	}{
		ivBase64,
		ciphertextBase64,
	})
}

func getTagItems(w http.ResponseWriter, r *http.Request) {
	// rate limit
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	rateLimit_key := ip + "get_tag_items"
	if !getRateLimiter(rateLimit_key, 30, 60*time.Second) { // 30 requests per minute
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		return
	}

	// parse request body
	var request struct {
		Tag struct{ 
			Type string `json:"type"`
			Id   string `json:"id"` } `json:"tag"`
		SessionId string `json:"session_id"`
	} // request body struct
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// get session key from session
	sessionKey, found := core.SessionKeyCache.Get(request.SessionId)
	if !found {
		http.Error(w, "Invalid session ID", http.StatusBadRequest)
		return
	}

	// get tag items data
	data, err := core.GetTagItems(request.Tag.Type, request.Tag.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// encrypt tag items data
	ivBase64, ciphertextBase64, err := crypto.EncryptAESGCM(sessionKey.([]byte), data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		IVBase64            string `json:"iv_base64"`
		CiphertextBase64    string `json:"ciphertext_base64"`
	}{
		ivBase64,
		ciphertextBase64,
	})
}

func getFilesMetaData(w http.ResponseWriter, r *http.Request) {
	// ip rate limit check
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	rateLimit_key := ip + "getFilesMetaData"
	if !getRateLimiter(rateLimit_key, 60, 60*time.Second) { // 60 requests per minute
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		return
	}

	// parse request body
	var folder struct {
		FileIds   []string `json:"file_ids"`
		SessionId string `json:"session_id"`
	} // request body struct
	err = json.NewDecoder(r.Body).Decode(&folder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// get session key from session
	sessionKey, found := core.SessionKeyCache.Get(folder.SessionId)
	if !found {
		http.Error(w, "Invalid session ID", http.StatusBadRequest)
		return
	}

	// get files metadata data
	data, err := core.GetFilesMetadata(&folder.FileIds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// encrypt files metadata data
	ivBase64, ciphertextBase64, err := crypto.EncryptAESGCM(sessionKey.([]byte), data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		IVBase64            string `json:"iv_base64"`
		CiphertextBase64    string `json:"ciphertext_base64"`
	}{
		ivBase64,
		ciphertextBase64,
	})
}

func getFoldersMetaData(w http.ResponseWriter, r *http.Request) {
	// ip rate limit check
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	rateLimit_key := ip + "getFoldersMetaData"
	if !getRateLimiter(rateLimit_key, 60, 60*time.Second) { // 60 requests per minute
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		return
	}

	// parse request body
	var folder struct {
		FolderIds []string `json:"folder_ids"`
		SessionId string `json:"session_id"`
	} // request body struct
	err = json.NewDecoder(r.Body).Decode(&folder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// get session key from session
	sessionKey, found := core.SessionKeyCache.Get(folder.SessionId)
	if !found {
		http.Error(w, "Invalid session ID", http.StatusBadRequest)
		return
	}

	// get folders metadata data
	data, err := core.GetFoldersMetadata(&folder.FolderIds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// encrypt folders metadata data
	ivBase64, ciphertextBase64, err := crypto.EncryptAESGCM(sessionKey.([]byte), data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		IVBase64            string `json:"iv_base64"`
		CiphertextBase64    string `json:"ciphertext_base64"`
	}{
		ivBase64,
		ciphertextBase64,
	})
}

func rename(w http.ResponseWriter, r *http.Request) {
	// rate limit
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	rateLimit_key := ip + "rename"
	if !getRateLimiter(rateLimit_key, 30, 60*time.Second) { // 30 requests per minute
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		return
	}

	// parse request body
	var request struct {
		IVBase64          string `json:"iv_base64"`
		CiphertextBase64  string `json:"ciphertext_base64"`
		SessionId string `json:"session_id"`
	} // request body struct
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// get session key from session
	sessionKey, found := core.SessionKeyCache.Get(request.SessionId)
	if !found {
		http.Error(w, "Invalid session ID", http.StatusBadRequest)
		return
	}

	// decrypt item data
	var requestItemData struct {
		Id   string `json:"id"`
		Type string `json:"type"`
		NewName string `json:"new_name"`
	}
	
	err = crypto.DecryptJSON(sessionKey.([]byte), request.IVBase64, request.CiphertextBase64, &requestItemData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// rename item
	err = core.Rename(requestItemData.Id, requestItemData.Type, requestItemData.NewName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}