package main

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

const yangFolder = "../uploads/"

func (s *srv) logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Print("\n")
		s.logger.Printf("REQUEST: %s %s %s", r.RemoteAddr, r.Method, r.URL)

		next.ServeHTTP(w, r)

		corHeader := "Access-Control-Allow-Origin"
		if w.Header().Get(corHeader) != "*" {
			w.Header().Set(corHeader, "*")
		}
	})
}

func connectionOk(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("[Success] backend active"))
}

func (s *srv) upload(w http.ResponseWriter, r *http.Request) {
	// Limit upload size (10 MB in this case)
	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("file")
	if err != nil {
		s.raiseError("[Error] retrieving zip file", err, w)
		return
	}

	defer file.Close()

	// Save the file locally
	zipPath := yangFolder + handler.Filename
	out, err := os.Create(zipPath)
	if err != nil {
		s.raiseError(fmt.Sprintf("[Error] creating (%s) zip file", handler.Filename), err, w)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		s.raiseError(fmt.Sprintf("[Error] copying (%s) zip file", handler.Filename), err, w)
		return
	}

	// Extract the zip file
	err = extractYangFolder(handler.Filename)
	if err != nil {
		s.raiseError(fmt.Sprintf("[Error] extracting (%s) yang files", handler.Filename), err, w)
		return
	}

	w.Write([]byte("[Success] repo uploaded"))
}

func extractYangFolder(filename string) error {
	zipPath := yangFolder + filename
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("[Error] reading ZIP file: %v", err)
	}
	defer r.Close()

	basename := strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))
	destFolder := filepath.Join(yangFolder + basename)

	yangFileCount := 0

	for _, f := range r.File {
		// Only extract .yang files
		if strings.HasSuffix(f.Name, ".yang") {
			yangFileCount += 1

			if yangFileCount == 1 {
				if _, err := os.Stat(destFolder); os.IsNotExist(err) {
					err := os.MkdirAll(destFolder, os.ModePerm)
					if err != nil {
						return fmt.Errorf("[Error] creating repo folder: %v", err)
					}
				}
			}

			fpath := filepath.Join(destFolder, filepath.Base(f.Name))

			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return fmt.Errorf("[Error] setting yang file write permission: %v", err)
			}
			defer outFile.Close()

			rc, err := f.Open()
			if err != nil {
				return fmt.Errorf("[Error] opening yang file into the repo: %v", err)
			}

			_, err = io.Copy(outFile, rc)
			rc.Close()

			if err != nil {
				return fmt.Errorf("[Error] copying (%s) yang file content: %v", f.Name, err)
			}
		}
	}

	err = os.Remove(zipPath)
	if err != nil {
		return fmt.Errorf("[Error] cleaning up (%s) zip file: %v", filename, err)
	}

	return nil
}

func (s *srv) uploadFile(w http.ResponseWriter, r *http.Request) {
	basename, ok := mux.Vars(r)["basename"]

	// Limit upload size (10 MB in this case)
	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("file")
	if err != nil {
		s.raiseError("[Error] retrieving .yang file", err, w)
		return
	}

	defer file.Close()

	folderPath := yangFolder
	if ok {
		folderPath = yangFolder + basename + "/"
	}

	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		s.raiseError(fmt.Sprintf("[Error] repo (%s) does not exist", basename), err, w)
		return
	}

	// Save the file locally
	filePath := folderPath + handler.Filename
	out, err := os.Create(filePath)
	if err != nil {
		s.raiseError(fmt.Sprintf("[Error] creating (%s) file", handler.Filename), err, w)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		s.raiseError(fmt.Sprintf("[Error] copying (%s) file", handler.Filename), err, w)
		return
	}

	w.Write([]byte("[Success] file uploaded"))
}

func (s *srv) list(w http.ResponseWriter, r *http.Request) {
	kind := mux.Vars(r)["kind"]

	type ListResponse struct {
		Name  string   `json:"name"`
		Files []string `json:"files,omitempty"`
	}

	var f []ListResponse

	if kind == "local" {
		dirEntires, err := os.ReadDir(yangFolder)
		if err != nil {
			s.raiseError("[Error] reading local yang repos", err, w)
			return
		}

		for _, entry := range dirEntires {
			if entry.IsDir() {
				folderName := entry.Name()

				repoEntires, err := os.ReadDir(yangFolder + folderName + "/")
				if err != nil {
					s.raiseError("[Error] reading local yang repos", err, w)
					return
				}

				var fEntry ListResponse
				fEntry.Name = folderName
				for _, entry := range repoEntires {
					if !entry.IsDir() {
						fEntry.Files = append(fEntry.Files, entry.Name())
					}
				}

				f = append(f, fEntry)
			}
		}

		var fEntry ListResponse
		fEntry.Name = ""
		for _, entry := range dirEntires {
			if !entry.IsDir() {
				yangFile := entry.Name()
				if strings.ToLower(filepath.Ext(yangFile)) == ".yang" {
					fEntry.Files = append(fEntry.Files, yangFile)
				}
			}
		}
		f = append(f, fEntry)

	} else if kind == "nsp" {
		intentTypeList, err := s.intentTypeSearch(0, 100)
		if err != nil {
			s.raiseError("[Error] fetching NSP Intent Types", err, w)

			err := s.revokeToken()
			if err != nil {
				s.logger.Printf("[Error] disconnecting with NSP (%s): %v", s.nsp.Ip, err)
				return
			}

			return
		}

		for _, entry := range intentTypeList {
			var fEntry ListResponse
			fEntry.Name = entry
			f = append(f, fEntry)
		}
	}

	b, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		s.raiseError("[Error] creating JSON", err, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func (s *srv) delete(w http.ResponseWriter, r *http.Request) {
	basename := mux.Vars(r)["basename"]
	folderPath := yangFolder + basename

	if _, err := os.Stat(folderPath); errors.Is(err, os.ErrNotExist) {
		s.raiseError(fmt.Sprintf("[Error] %s repo does not exist", basename), err, w)
		return
	}

	err := os.RemoveAll(folderPath)
	if err != nil {
		s.raiseError("[Error] during repo deletion", err, w)
		return
	}

	w.Write([]byte("[Success] Local repo (" + basename + ") deleted"))
}

func (s *srv) deleteFile(w http.ResponseWriter, r *http.Request) {
	basename, ok := mux.Vars(r)["basename"]
	yangFile := mux.Vars(r)["yang"]

	filePath := yangFolder + yangFile
	if ok {
		filePath = yangFolder + basename + "/" + yangFile
	}

	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		s.raiseError(fmt.Sprintf("[Error] %s repo does not exist", basename), err, w)
		return
	}

	err := os.RemoveAll(filePath)
	if err != nil {
		s.raiseError("[Error] during file deletion", err, w)
		return
	}

	w.Write([]byte(fmt.Sprintf("[Success] file (%s/%s) deleted", basename, yangFile)))
}

func (s *srv) nspConnect(w http.ResponseWriter, r *http.Request) {
	err := json.NewDecoder(r.Body).Decode(&s.nsp)
	if err != nil {
		s.raiseError("[Error] decoding NSP connect request", err, w)
		return
	}

	if s.nsp.Ip == "" || s.nsp.User == "" || s.nsp.Pass == "" {
		s.raiseError("[Error] NSP credentials are missing", nil, w)
		return
	}

	err = s.getToken()
	if err != nil {
		s.raiseError("[Error] making NSP connection", err, w)
		return
	}

	go func() {
		for {
			s.Lock()
			if s.nsp.token.ExpiresIn == 0 {
				s.Unlock()
				return
			}

			timeLeft := s.nsp.token.ExpiresIn
			time.Sleep(1 * time.Second)
			timeLeft--

			if timeLeft <= 10 {
				s.logger.Println("[Info] NSP Access renewal initiated")
				err := s.revokeToken()
				if err != nil {
					s.logger.Printf("[Error] disconnecting with NSP (%s): %v", s.nsp.Ip, err)
					s.Unlock()
					return
				}
				err = s.getToken()
				if err != nil {
					s.logger.Printf("[Error] reconnecting with NSP (%s): %v", s.nsp.Ip, err)
					s.Unlock()
					return
				}
				s.logger.Println("[Success] NSP Access renewed")
			} else {
				s.nsp.token.ExpiresIn = timeLeft
				//s.logger.Printf("[Success] NSP Access expires in %d", timeLeft)
			}

			s.Unlock()
		}
	}()

	w.Write([]byte("[Success] NSP connected"))
}

func (s *srv) nspDisconnect(w http.ResponseWriter, r *http.Request) {
	if s.nsp.Ip == "" {
		s.raiseError("[Error] NSP is not connected", nil, w)
		return
	}
	err := s.revokeToken()
	if err != nil {
		s.raiseError(fmt.Sprintf("[Error] disconnecting with NSP (%s)", s.nsp.Ip), err, w)
		return
	}

	w.Write([]byte("[Success] NSP disconnected"))
}

func (s *srv) nspIsConnected(w http.ResponseWriter, r *http.Request) {
	if s.nsp.token.AccessToken == "" {
		s.raiseError("[Error] NSP not connected", nil, w)
		return
	}

	nspExport := NspAccessExport{
		Ip:   s.nsp.Ip,
		User: s.nsp.User,
	}

	b, err := json.MarshalIndent(nspExport, "", "  ")
	if err != nil {
		s.raiseError("Error creating JSON", err, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}
