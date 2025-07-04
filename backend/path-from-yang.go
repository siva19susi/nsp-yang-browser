package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/gorilla/mux"
	"github.com/openconfig/goyang/pkg/yang"
)

type Dependency struct {
	Lso        []string `json:"lso"`
	Module     []string `json:"Module"`
	IntentType []string `json:"intent-type"`
}

const commonFolder = "../common/"

var (
	dependency = Dependency{
		Lso: []string{
			"nsp-lso-manager.yang",
			"nsp-lso-operation.yang",
			"nsp-model-extensions.yang",
			"ietf-yang-types.yang",
			"ietf-inet-types.yang",
		},
		Module: []string{
			"nsp-model-extensions.yang",
		},
		IntentType: []string{
			"ietf-inet-types.yang",
			"ietf-yang-types.yang",
			"webfwk-ui-metadata.yang",
		},
	}
)

// Handler for generating schema from YANG files
func (s *srv) pathFromYang(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]

	saveParam := r.URL.Query().Get("save")
	save := saveParam == "true"

	pathSegments := strings.Split(r.URL.Path, "/")
	kind := pathSegments[1]
	subKind := pathSegments[2]

	app := &App{
		SchemaTree: &yang.Entry{
			Dir: make(map[string]*yang.Entry),
		},
		modules: yang.NewModules(),
	}

	var definitions []YangDefinition

	switch kind {
	case "offline":
		var files []string
		dirPath := filepath.Join(commonFolder, name)
		dirFiles, err := os.ReadDir(dirPath)
		if err != nil {
			s.raiseError("directory missing", err, w)
			return
		}
		for _, file := range dirFiles {
			files = append(files, filepath.Join(dirPath, file.Name()))
		}

		if err := app.readYangModules(files); err != nil {
			s.raiseError("error generating YANG schema", err, w)
			return
		}
	case "nsp":
		switch subKind {
		case "intent-type":
			{
				yangModules, err := s.intentTypeYangModules(name)
				if err != nil {
					s.raiseError("", err, w)
					return
				}
				for _, yangModule := range yangModules {
					definitions = append(definitions, YangDefinition{
						Name:       yangModule.Name,
						Definition: yangModule.YangContent,
					})
				}

				definitions, err = loadDependencyDefinition(definitions, dependency.IntentType)
				if err != nil {
					s.raiseError("", err, w)
					return
				}

				if err := app.definitionToSchema(definitions); err != nil {
					s.raiseError("error generating YANG schema", err, w)
					return
				}
			}
		case "lso-operation":
			{
				operationName, operationYang, err := s.getLsoOperationModel(name)
				if err != nil {
					s.raiseError("", err, w)
					return
				}
				definitions = append(definitions, YangDefinition{
					Name:       operationName,
					Definition: operationYang,
				})

				definitions, err = loadDependencyDefinition(definitions, dependency.Lso)
				if err != nil {
					s.raiseError("", err, w)
					return
				}

				if err := app.definitionToSchema(definitions); err != nil {
					s.raiseError("error generating YANG schema", err, w)
					return
				}
			}
		}
	}

	result, err := app.pathCmdRun()
	if err != nil {
		s.raiseError("[Error] running path command", err, w)
		return
	}

	if save {
		if err := s.saveJsonFile(subKind, name, result); err != nil {
			s.raiseError("", err, w)
			return
		}
		writeResponse(w, "success", fmt.Sprintf("%s/%s.json was saved", subKind, name))
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write(result)
	}
}

func loadDependencyDefinition(definitions []YangDefinition, dependents []string) ([]YangDefinition, error) {
	dirPath := filepath.Join(commonFolder)
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return definitions, fmt.Errorf("common yang folder missing %v", err)
	}

	defintionsAlreadyContains := func(name string) bool {
		for _, def := range definitions {
			if def.Name == name {
				return true
			}
		}
		return false
	}

	for _, file := range files {
		fileName := file.Name()

		if !file.IsDir() && slices.Contains(dependents, fileName) && !defintionsAlreadyContains(fileName) {
			filePath := filepath.Join(commonFolder, fileName)
			content, err := os.ReadFile(filePath)
			if err != nil {
				return definitions, fmt.Errorf("error reading %s", fileName)
			}
			definitions = append(definitions, YangDefinition{
				Name:       fileName,
				Definition: string(content),
			})
		}
	}

	return definitions, nil
}
