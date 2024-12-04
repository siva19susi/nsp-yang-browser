package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
	"github.com/openconfig/goyang/pkg/yang"
)

// Add YANG file to the list if not already present
func addYangFileIfMissing(files []string, yangFileName string) []string {
	for _, file := range files {
		if strings.Contains(file, yangFileName) {
			return files
		}
	}
	return append(files, filepath.Join(yangFolder, yangFileName))
}

// Load YANG file content into IntentTypeYangModule
func loadYangModule(yangFileName string) (IntentTypeYangModule, error) {
	filePath := filepath.Join(yangFolder, yangFileName)
	content, err := os.ReadFile(filePath)
	if err != nil {
		return IntentTypeYangModule{}, fmt.Errorf("[Error] reading %s file: %v", yangFileName, err)
	}
	return IntentTypeYangModule{
		Name:        yangFileName,
		YangContent: string(content),
	}, nil
}

// Get all common YANG files in the folder
func getCommonYangFiles() ([]string, error) {
	files, err := os.ReadDir(yangFolder)
	if err != nil {
		return nil, fmt.Errorf("[Error] reading yang repo: %v", err)
	}

	var commonYangs []string
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".yang" {
			commonYangs = append(commonYangs, file.Name())
		}
	}
	return commonYangs, nil
}

// Read and parse YANG files in the specified directory
func readYangFilesFromDir(dirPath string, commonYangs []string) ([]string, error) {
	dirEntries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("[Error] reading yang repo: %v", err)
	}

	var files []string
	for _, entry := range dirEntries {
		files = append(files, filepath.Join(dirPath, entry.Name()))
	}

	for _, commonYang := range commonYangs {
		files = addYangFileIfMissing(files, commonYang)
	}

	return files, nil
}

// Handler for generating schema from YANG files
func (s *srv) pathFromYang(w http.ResponseWriter, r *http.Request) {
	kind := mux.Vars(r)["kind"]
	basename := mux.Vars(r)["basename"]

	app := &App{
		SchemaTree: &yang.Entry{
			Dir: make(map[string]*yang.Entry),
		},
		modules: yang.NewModules(),
	}

	commonYangs, err := getCommonYangFiles()
	if err != nil {
		s.raiseError("[Error] reading common YANG files", err, w)
		return
	}

	switch kind {
	case "local":
		dirPath := filepath.Join(yangFolder, basename)
		files, err := readYangFilesFromDir(dirPath, commonYangs)
		if err != nil {
			s.raiseError("[Error] preparing YANG files", err, w)
			return
		}

		if err := app.readYangModules(files); err != nil {
			s.raiseError("[Error] generating YANG schema", err, w)
			return
		}

	case "nsp":
		yangModules, err := s.intentTypeYangModules(basename)
		if err != nil {
			s.raiseError("[Error] fetching YANG modules", err, w)
			return
		}

		for _, commonYang := range commonYangs {
			module, err := loadYangModule(commonYang)
			if err != nil {
				s.raiseError("[Error] loading common YANG module", err, w)
				return
			}
			yangModules = append(yangModules, module)
		}

		var definitions []YangDefinition
		for _, module := range yangModules {
			definitions = append(definitions, YangDefinition{
				Name:       module.Name,
				Definition: module.YangContent,
			})
		}

		if err := app.definitionToSchema(definitions); err != nil {
			s.raiseError("[Error] generating YANG schema", err, w)
			return
		}
	}

	result, err := app.pathCmdRun()
	if err != nil {
		s.raiseError("[Error] running path command", err, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
}
