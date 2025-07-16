package httpserver

import "net/http"

func RegisterRoutes(mux *http.ServeMux) {
	// mux.Handle("/", http.FileServer(http.Dir("./client")))
	// mux.HandleFunc("/rootDirPath", getRootDirPath)
	// mux.HandleFunc("/folder", getFolder)
	// mux.HandleFunc("/files", getFile)
	
	mux.HandleFunc("POST /handshake/connect", handShakeConnect)
	mux.HandleFunc("POST /handshake/authenticate", handShakeAuthenticate)
	mux.HandleFunc("POST /refreshSessionKey", refreshSessionKey)
	mux.HandleFunc("POST /getFolder", getFolder)
	mux.HandleFunc("POST /getTags", getTags)
	mux.HandleFunc("POST /getTagItems", getTagItems)
	mux.HandleFunc("POST /getFilesMetaData", getFilesMetaData)
}