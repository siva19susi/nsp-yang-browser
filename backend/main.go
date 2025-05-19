package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

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
	s.be.HandleFunc("/upload/file/{name}", s.uploadFile).Methods("POST")
	s.be.HandleFunc("/uploaded/all", s.uploadedAll).Methods("GET")
	s.be.HandleFunc("/uploaded/{name}", s.uploadedSpecific).Methods("GET")
	s.be.HandleFunc("/uploaded/{name}/paths", s.pathFromYang).Methods("GET")

	// DOWNLOAD HANDLER
	s.be.HandleFunc("/download/{name}", s.downloadBundle).Methods("GET")
	s.be.HandleFunc("/download/{name}/file/{yang}", s.downloadYang).Methods("GET")

	// DELETE HANDLERS
	s.be.HandleFunc("/delete/{name}", s.delete).Methods("DELETE")
	s.be.HandleFunc("/delete/{name}/file/{yang}", s.deleteFile).Methods("DELETE")
	s.be.HandleFunc("/delete/file/{yang}", s.deleteFile).Methods("DELETE")

	// NSP HANDLERS
	s.be.HandleFunc("/nsp/connect", s.nspConnect).Methods("POST")
	s.be.HandleFunc("/nsp/disconnect", s.nspDisconnect).Methods("POST")
	s.be.HandleFunc("/nsp/isConnected", s.nspIsConnected).Methods("GET")

	s.be.HandleFunc("/nsp/modules", s.getNspModules).Methods("GET")
	s.be.HandleFunc("/nsp/module/{name}/paths", s.getNspModulePaths).Methods("GET")

	s.be.HandleFunc("/nsp/intent-types", s.getIntentTypes).Methods("GET")
	s.be.HandleFunc("/nsp/intent-type/{name}/paths", s.pathFromYang).Methods("GET")
	s.be.HandleFunc("/nsp/intent-type/{name}/intents", s.getIntents).Methods("GET")

	s.be.HandleFunc("/nsp/find", s.nspFind).Methods("POST")
	s.be.HandleFunc("/nsp/intent-explorer", s.intentExplorer).Methods("POST")

	s.be.HandleFunc("/nsp/lso-operations", s.getLsoOperations).Methods("GET")
	s.be.HandleFunc("/nsp/lso-operation/{name}/paths", s.pathFromYang).Methods("GET")

	s.be.HandleFunc("/nsp/telemetry-types", s.getTelemetryTypes).Methods("GET")
	s.be.HandleFunc("/nsp/telemetry-type/definition", s.getTelemetryTypeDefinition).Methods("POST")

	muxSrv := &http.Server{
		Addr:    ":8080",
		Handler: s.be,
	}

	// HANDLE CTRL+C
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		// START HTTP SERVER
		s.logger.Printf("Access API with baseURL - http://localhost:8080")
		if err := muxSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Fatalf("ListenAndServe Error: %v", err)
		}
	}()

	// WAIT for Ctrl+C SIGNAL
	<-stop

	// CREATE A DEADLINE TO WAIT FOR SHUTDOWN
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// GRACEFULLY REVOKING NSP TOKEN
	if s.nsp.token.AccessToken != "" {
		s.revokeToken()
	}

	// ATTEMPT GRACEFUL SHUTDOWN
	if err := muxSrv.Shutdown(ctx); err != nil {
		s.logger.Fatalf("Server forced to shutdown: %v\n", err)
	}
}
