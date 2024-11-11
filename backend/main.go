package main

import (
	"fmt"
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

type NspAccessExport struct {
	Ip   string `json:"ip"`
	User string `json:"user"`
}

type srv struct {
	sync.Mutex
	be     *mux.Router
	logger *log.Logger
	nsp    NspAccess
}

func (s *srv) raiseError(mess string, err error, w http.ResponseWriter) {
	errMessage := fmt.Sprintf("%s: %v", mess, err)
	s.logger.Println(errMessage)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	http.Error(w, errMessage, http.StatusInternalServerError)
}

func main() {
	s := &srv{
		be:     mux.NewRouter(),
		logger: log.New(os.Stderr, "", log.LstdFlags),
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

	s.be.Use(s.logMiddleware)
	s.be.HandleFunc("/", connectionOk).Methods("GET")
	s.be.HandleFunc("/upload", s.upload).Methods("POST")
	s.be.HandleFunc("/upload/{basename}", s.uploadFile).Methods("POST")
	s.be.HandleFunc("/list/{kind}", s.list).Methods("GET")
	s.be.HandleFunc("/delete/{basename}", s.delete).Methods("DELETE")
	s.be.HandleFunc("/delete/{basename}/file/{yang}", s.deleteFile).Methods("DELETE")
	s.be.HandleFunc("/generate/{kind}/{basename}", s.pathCmdRun).Methods("GET")
	s.be.HandleFunc("/nsp/connect", s.nspConnect).Methods("POST")
	s.be.HandleFunc("/nsp/disconnect", s.nspDisconnect).Methods("POST")
	s.be.HandleFunc("/nsp/isConnected", s.nspIsConnected).Methods("GET")

	s.logger.Printf("Access API with baseURL - http://localhost:8080")
	err := http.ListenAndServe(":8080", s.be)
	if err != nil {
		s.logger.Printf("ListenAndServer Error: %s", err.Error())
	}
}
