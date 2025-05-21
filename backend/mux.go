package main

import (
	"archive/zip"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/openconfig/goyang/pkg/yang"
)

const yangFolder = "../offline/"

func (s *srv) logMiddleware(next http.Handler) http.Handler {
	const corHeader = "Access-Control-Allow-Origin"

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		s.logger.Printf("REQUEST: %s %s %s", r.RemoteAddr, r.Method, r.URL)

		defer func() {
			s.logger.Printf("RESPONSE: %s %s %s completed in %v", r.RemoteAddr, r.Method, r.URL, time.Since(start))
		}()

		// Set CORS header
		w.Header().Set(corHeader, "*")

		next.ServeHTTP(w, r)
	})
}

// WRITE RESPONSE PLAINTEXT
func writeResponse(w http.ResponseWriter, status string, msg string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if status == "error" && msg != "" {
		http.Error(w, msg, http.StatusInternalServerError)
	} else if status == "success" {
		w.Write([]byte(msg))
	} else {
		http.Error(w, "unknown error", http.StatusInternalServerError)
	}
}

// WRITE RESPONSE JSON
func writeJsonResponse(w http.ResponseWriter, b []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

// RAISE ERRORS
func (s *srv) raiseError(msg string, err error, w http.ResponseWriter) {
	if err != nil {
		s.logger.Printf(msg, err)
		writeResponse(w, "error", fmt.Sprintf("%s / %v", msg, err))
	} else {
		writeResponse(w, "error", msg)
	}
}

// Helper function to handle file creation and writing
func saveFile(file io.Reader, filepath string) error {
	out, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("file creation failed: %v", err)
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		return fmt.Errorf("copying file content failed: %v", err)
	}

	return nil
}

func (s *srv) saveJsonFile(moduleName string, fileName string, content []byte) error {
	u := uuid.New()
	uuid := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(u[:])
	saveFileName := fmt.Sprintf("%s__%s__%s__%s__%s.json", uuid, s.nsp.Ip, s.nsp.token.ConnectTime, moduleName, fileName)
	saveFilePath := filepath.Join(yangFolder, saveFileName)
	err := os.WriteFile(saveFilePath, content, 0644)
	if err != nil {
		return fmt.Errorf("error saving to local directory %v", err)
	}
	return nil
}

// BACKEND CONNECTION VERIFICATION
func connectionOk(w http.ResponseWriter, r *http.Request) {
	writeResponse(w, "success", "Backend active")
}

// UPLOAD YANG REPO ZIP
func (s *srv) upload(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("file")
	if err != nil {
		s.raiseError("retrieving zip file failed", err, w)
		return
	}
	defer file.Close()

	zipPath := yangFolder + handler.Filename
	if err := saveFile(file, zipPath); err != nil {
		s.raiseError(fmt.Sprintf("saving zip file failed %s", handler.Filename), err, w)
		return
	}

	if err := s.extractYangFolder(handler.Filename); err != nil {
		s.raiseError(fmt.Sprintf("extracting yang files from %s failed", handler.Filename), err, w)
		return
	}

	writeResponse(w, "success", "Repo uploaded")
}

// UNZIP YANG REPO
func (s *srv) extractYangFolder(filename string) error {
	zipPath := yangFolder + filename
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("reading zip file failed: %v", err)
	}
	defer r.Close()

	basename := strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))
	destFolder := filepath.Join(yangFolder + basename)

	if err := os.MkdirAll(destFolder, os.ModePerm); err != nil {
		return fmt.Errorf("creating repo folder failed: %v", err)
	}

	var files []string

	yangFileCount := 0
	for _, f := range r.File {
		if strings.HasSuffix(f.Name, ".yang") {
			yangFileCount++
			fpath := filepath.Join(destFolder, filepath.Base(f.Name))
			files = append(files, fpath)

			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return fmt.Errorf("setting yang file write permission failed: %v", err)
			}
			defer outFile.Close()

			rc, err := f.Open()
			if err != nil {
				return fmt.Errorf("opening yang file failed: %v", err)
			}
			defer rc.Close()

			_, err = io.Copy(outFile, rc)
			if err != nil {
				return fmt.Errorf("copying file content failed: %v", err)
			}
		}
	}

	if err = os.Remove(zipPath); err != nil {
		return err
	}

	var failed = true
	defer func() {
		if failed {
			os.RemoveAll(destFolder)
		}
	}()

	app := &App{
		SchemaTree: &yang.Entry{
			Dir: make(map[string]*yang.Entry),
		},
		modules: yang.NewModules(),
	}

	if err := app.readYangModules(files); err != nil {
		return fmt.Errorf("error generating YANG schema %v", err)
	}

	result, err := app.pathCmdRun()
	if err != nil {
		return fmt.Errorf("error running path command %v", err)
	}

	u := uuid.New()
	uuid := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(u[:])
	currentTime := time.Now().Format("2006-01-02-15-04-05")
	fileName := fmt.Sprintf("%s__%s__%s__%s__%s.json", uuid, "local", currentTime, "uploaded", basename)
	saveFile := filepath.Join(yangFolder, fileName)
	err = os.WriteFile(saveFile, result, 0644)
	if err != nil {
		return fmt.Errorf("error saving to local directory %v", err)
	}

	failed = false
	return os.RemoveAll(destFolder)
}

// LIST KIND
func (s *srv) uploadedAll(w http.ResponseWriter, r *http.Request) {
	var f []string

	entries, err := os.ReadDir(yangFolder)
	if err != nil {
		s.raiseError("error reading offline folder", err, w)
		return
	}

	for _, entry := range entries {
		name := strings.Replace(entry.Name(), ".json", "", -1)
		if !entry.IsDir() && !strings.HasPrefix(name, ".") {
			f = append(f, name)
		}
	}

	b, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		s.raiseError("JSON creation failed", err, w)
		return
	}

	writeJsonResponse(w, b)
}

func (s *srv) uploadedInfo(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	entries, err := os.ReadDir(yangFolder)
	if err != nil {
		s.raiseError("error reading offline folder", err, w)
		return
	}

	type Uploaded struct {
		Id        string `json:"id"`
		NspIp     string `json:"nspIp"`
		Timestamp string `json:"timestamp"`
		Module    string `json:"module"`
		Name      string `json:"name"`
	}

	var info Uploaded
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasPrefix(entry.Name(), id) {
			parts := strings.Split(entry.Name(), "__")
			info = Uploaded{
				Id:        parts[0],
				NspIp:     parts[1],
				Timestamp: parts[2],
				Module:    parts[3],
				Name:      strings.Replace(parts[4], ".json", "", -1),
			}
			break
		}
	}

	response, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		s.raiseError("JSON creation failed", err, w)
		return
	}

	writeJsonResponse(w, response)
}

func (s *srv) uploadedPaths(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	entries, err := os.ReadDir(yangFolder)
	if err != nil {
		s.raiseError("error reading offline folder", err, w)
		return
	}

	var filename string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasPrefix(entry.Name(), id) {
			filename = entry.Name()
			break
		}
	}

	data, err := os.ReadFile(filepath.Join(yangFolder, filename))
	if err != nil {
		s.raiseError("Unable to read file", err, w)
		return
	}

	writeJsonResponse(w, data)
}

func (s *srv) downloadBundle(w http.ResponseWriter, r *http.Request) {
	folderName := mux.Vars(r)["name"]
	folderPath := yangFolder + folderName
	zipFileName := folderName + ".zip"

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename=\""+zipFileName+"\"")

	zipWriter := zip.NewWriter(w)
	defer zipWriter.Close()

	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		relPath := strings.TrimPrefix(path, folderPath)
		relPath = strings.TrimLeft(relPath, string(filepath.Separator))

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name = relPath
		header.Method = zip.Deflate

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		_, err = io.Copy(writer, file)
		return err
	})

	if err != nil {
		http.Error(w, "Failed to create zip: "+err.Error(), http.StatusInternalServerError)
	}
}

// DELETE FOLDER OR REPO
func (s *srv) delete(w http.ResponseWriter, r *http.Request) {
	basename := mux.Vars(r)["name"]
	folderPath := yangFolder + basename

	if _, err := os.Stat(folderPath); errors.Is(err, os.ErrNotExist) {
		s.raiseError(fmt.Sprintf("%s repo does not exist", basename), err, w)
		return
	}

	if err := os.RemoveAll(folderPath); err != nil {
		s.raiseError("error during repo deletion", err, w)
		return
	}

	writeResponse(w, "success", fmt.Sprintf("Local repo (%s) deleted", basename))
}

// NSP CONNECT
func (s *srv) nspConnect(w http.ResponseWriter, r *http.Request) {
	if err := json.NewDecoder(r.Body).Decode(&s.nsp); err != nil {
		s.raiseError("decoding NSP connect request failed", err, w)
		return
	}

	if s.nsp.Ip == "" || s.nsp.User == "" || s.nsp.Pass == "" {
		s.raiseError("NSP credentials are missing", nil, w)
		return
	}

	if err := s.getToken(); err != nil {
		s.raiseError("error making NSP connection", err, w)
		return
	}

	go s.tokenRefreshRoutine()

	writeResponse(w, "success", "NSP connected")
}

// Token refresh routine
func (s *srv) tokenRefreshRoutine() {
	for {
		tokenExpiresIn := s.nsp.token.ExpiresIn

		if tokenExpiresIn == 0 {
			return
		}

		waitTime := max(time.Duration(s.nsp.token.ExpiresIn-15)*time.Second, 0)
		time.Sleep(waitTime)

		s.logger.Println("[Info] NSP Access renewal initiated")

		if err := s.getToken(); err != nil {
			s.logger.Printf("reconnecting to NSP (%s) failed: %v", s.nsp.Ip, err)
			return
		}
		s.logger.Println("[Info] NSP Access renewed")
	}
}

// NSP DISCONNECT
func (s *srv) nspDisconnect(w http.ResponseWriter, r *http.Request) {
	if s.nsp.Ip == "" {
		s.raiseError("NSP is not connected", nil, w)
		return
	}

	if err := s.revokeToken(); err != nil {
		s.raiseError(fmt.Sprintf("disconnecting from NSP (%s) failed", s.nsp.Ip), err, w)
		return
	}

	writeResponse(w, "success", "NSP disconnected")
}

// NSP IS CONNECTED
func (s *srv) nspIsConnected(w http.ResponseWriter, r *http.Request) {
	if s.nsp.Ip == "" {
		s.raiseError("NSP is not connected", nil, w)
		return
	}

	if err := s.getNspStatus(); err != nil {
		s.raiseError("error getting health info", err, w)
		s.nspReset()
		return
	}

	type NspAccessExport struct {
		Ip   string `json:"ip"`
		User string `json:"user"`
	}

	nspExport := NspAccessExport{
		Ip:   s.nsp.Ip,
		User: s.nsp.User,
	}

	response, err := json.MarshalIndent(nspExport, "", "  ")
	if err != nil {
		s.raiseError("JSON creation failed", err, w)
		return
	}

	writeJsonResponse(w, response)
}

// GET NSP MODULES
func (s *srv) getNspModules(w http.ResponseWriter, r *http.Request) {
	modules, err := s.fetchModules()
	if err != nil {
		s.raiseError("error fetching YANG modules", err, w)
		return
	}

	response, err := json.MarshalIndent(modules, "", "  ")
	if err != nil {
		s.raiseError("error creating JSON", err, w)
		return
	}
	writeJsonResponse(w, response)
}

// GET NSP MODULE PATHS
func (s *srv) getNspModulePaths(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]

	saveParam := r.URL.Query().Get("save")
	save := saveParam == "true"

	response, err := s.fetchYangDefinition(name)
	if err != nil {
		s.raiseError("error fetching YANG module definition", err, w)
		return
	}

	if save {
		if err := s.saveJsonFile("module", name, response); err != nil {
			s.raiseError("", err, w)
			return
		}
		writeResponse(w, "success", fmt.Sprintf("module/%s.json was saved", name))
	} else {
		writeJsonResponse(w, response)
	}
}

// GET NSP INTENT TYPES
func (s *srv) getIntentTypes(w http.ResponseWriter, r *http.Request) {
	type SearchResponse struct {
		Total       int      `json:"total"`
		PageCount   int      `json:"pageCount"`
		IntentTypes []string `json:"intentTypes,omitempty"`
	}

	var getUrlQuery = func(r *http.Request, key string, defaultVal int) int {
		valStr := r.URL.Query().Get(key)
		if valStr == "" {
			return defaultVal
		}
		val, err := strconv.Atoi(valStr)
		if err != nil {
			return defaultVal
		}
		return val
	}

	limit := getUrlQuery(r, "limit", 30)
	page := getUrlQuery(r, "page", 1)

	filter := r.URL.Query().Get("filter")

	intentTypes, total, err := s.intentTypeSearch(page-1, limit, filter)
	if err != nil {
		s.raiseError("fetching NSP intent types failed", err, w)
		return
	}

	searched := SearchResponse{
		Total:       total,
		PageCount:   len(intentTypes),
		IntentTypes: intentTypes,
	}

	b, err := json.MarshalIndent(searched, "", "  ")
	if err != nil {
		s.raiseError("JSON creation failed", err, w)
		return
	}

	writeJsonResponse(w, b)
}

// GET NSP INTENT TYPES
func (s *srv) getIntents(w http.ResponseWriter, r *http.Request) {
	intentType := mux.Vars(r)["name"]

	lastInd := strings.LastIndex(intentType, "_")
	if lastInd == -1 {
		s.raiseError(fmt.Sprintf("Invalid intent type format: %s", intentType), nil, w)
		return
	}

	name := intentType[:lastInd]
	version := intentType[lastInd+2:]

	targets, err := s.getIntentTargets(name, version, 0)
	if err != nil {
		s.raiseError("fetching NSP intent types failed", err, w)
		return
	}

	b, err := json.MarshalIndent(targets, "", "  ")
	if err != nil {
		s.raiseError("JSON creation failed", err, w)
		return
	}

	writeJsonResponse(w, b)
}
