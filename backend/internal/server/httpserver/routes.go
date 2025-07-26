package httpserver

import "net/http"

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /handshake/connect", handShakeConnect)
	mux.HandleFunc("POST /handshake/authenticate", handShakeAuthenticate)
	mux.HandleFunc("POST /refreshSessionKey", refreshSessionKey)
	mux.HandleFunc("GET /getRootDirId", getRootDirId)
	mux.HandleFunc("POST /getFolder", getFolder)
	mux.HandleFunc("POST /getTags", getTags)
	mux.HandleFunc("POST /getTagItems", getTagItems)
	mux.HandleFunc("POST /getFilesMetaData", getFilesMetaData)
	mux.HandleFunc("POST /getFoldersMetaData", getFoldersMetaData)
	mux.HandleFunc("POST /rename", rename)
}