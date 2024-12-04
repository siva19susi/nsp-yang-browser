package main

import (
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/mux"
)

type TokenDetail struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

type NspAccess struct {
	Ip    string `json:"ip"`
	User  string `json:"user"`
	Pass  string `json:"pass"`
	token TokenDetail
}

type srv struct {
	sync.Mutex
	be     *mux.Router
	logger *log.Logger
	nsp    NspAccess
}

func main() {
	s := &srv{
		be:     mux.NewRouter(),
		logger: log.New(os.Stderr, "\n[BACKEND] ", log.LstdFlags),
		nsp: NspAccess{
			Ip:   "",
			User: "",
			Pass: "",
			token: TokenDetail{
				AccessToken:  "",
				RefreshToken: "",
				TokenType:    "",
				ExpiresIn:    0,
			},
		},
	}

	// ENSURE YANG FOLDER EXISTS
	if err := os.MkdirAll(yangFolder, os.ModePerm); err != nil {
		log.Fatalf("creating yang folder failed: %v", err)
	}

	// MIDDLEWARE
	s.be.Use(s.logMiddleware)

	// BASE API HEALTH CHECK
	s.be.HandleFunc("/", connectionOk).Methods("GET")

	// UPLOAD HANDLERS
	s.be.HandleFunc("/upload", s.upload).Methods("POST")
	s.be.HandleFunc("/upload/file", s.uploadFile).Methods("POST")
	s.be.HandleFunc("/upload/file/{basename}", s.uploadFile).Methods("POST")

	// LIST HANDLERS
	s.be.HandleFunc("/list/{kind}", s.list).Methods("GET")

	// DELETE HANDLERS
	s.be.HandleFunc("/delete/{basename}", s.delete).Methods("DELETE")
	s.be.HandleFunc("/delete/{basename}/file/{yang}", s.deleteFile).Methods("DELETE")
	s.be.HandleFunc("/delete/file/{yang}", s.deleteFile).Methods("DELETE")

	// GENERATE PATH HANDLER
	s.be.HandleFunc("/generate/{kind}/{basename}", s.pathFromYang).Methods("GET")

	// NSP HANDLERS
	s.be.HandleFunc("/nsp/connect", s.nspConnect).Methods("POST")
	s.be.HandleFunc("/nsp/disconnect", s.nspDisconnect).Methods("POST")
	s.be.HandleFunc("/nsp/isConnected", s.nspIsConnected).Methods("GET")
	s.be.HandleFunc("/nsp/modules", s.getNspModules).Methods("GET")
	s.be.HandleFunc("/nsp/module/{module}", s.getNspModules).Methods("GET")

	// START HTTP SERVER
	s.logger.Printf("Access API with baseURL - http://localhost:8080")
	if err := http.ListenAndServe(":8080", s.be); err != nil {
		s.logger.Fatalf("ListenAndServe Error: %v", err)
	}
}
