package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/openconfig/goyang/pkg/yang"
)

type YangDefinition struct {
	Name       string `json:"module-name"`
	Definition string `json:"yang-definition"`
}

type ModuleSet struct {
	Output struct {
		Modules []string `json:"module-set"`
	} `json:"nsp-yang-modules:output"`
}

type ModuleDefinition struct {
	Output struct {
		Result []YangDefinition `json:"result"`
	} `json:"nsp-yang-modules:output"`
}

// fetchModules retrieves the list of available NSP modules
func (s *srv) fetchModules() ([]string, error) {
	url := fmt.Sprintf("https://%s/restconf/operations/nsp-yang-modules:get-yang-module-sets", s.nsp.Ip)

	resp, err := s.makeHTTPRequest("POST", url, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("[Error] fetching NSP modules: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("[Error] fetching NSP modules, status: %d", resp.StatusCode)
	}

	var modules ModuleSet
	if err := json.NewDecoder(resp.Body).Decode(&modules); err != nil {
		return nil, fmt.Errorf("[Error] decoding modules response: %v", err)
	}

	// IETF module has augment field which go-yang does not support
	targetModules := []string{}
	for _, m := range modules.Output.Modules {
		if m != "ietf" {
			targetModules = append(targetModules, m)
		}
	}

	return targetModules, nil
}

// fetchYangDefinition retrieves YANG definitions for a specified module
func (s *srv) fetchYangDefinition(module string) ([]byte, error) {
	payload := strings.NewReader(fmt.Sprintf(`{"input": {"module-set-name": "%s"}}`, module))
	url := fmt.Sprintf("https://%s/restconf/operations/nsp-yang-modules:get-yang-modules-definitions", s.nsp.Ip)

	resp, err := s.makeHTTPRequest("POST", url, payload, nil)
	if err != nil {
		return nil, fmt.Errorf("[Error] fetching YANG definition: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("[Error] fetching YANG definition, status: %d", resp.StatusCode)
	}

	var yangDefs ModuleDefinition
	if err := json.NewDecoder(resp.Body).Decode(&yangDefs); err != nil {
		return nil, fmt.Errorf("[Error] decoding YANG definition response: %v", err)
	}

	yangDefinitions := filterValidYangDefinitions(yangDefs.Output.Result)

	app := &App{
		SchemaTree: &yang.Entry{
			Dir: make(map[string]*yang.Entry),
		},
		modules: yang.NewModules(),
	}

	if err := app.definitionToSchema(yangDefinitions); err != nil {
		return nil, fmt.Errorf("[Error] converting definitions to schema: %v", err)
	}

	return app.pathCmdRun()
}

// filterValidYangDefinitions filters out YANG definitions with empty names
func filterValidYangDefinitions(definitions []YangDefinition) []YangDefinition {
	var validDefinitions []YangDefinition
	for _, def := range definitions {
		if def.Name != "" {
			validDefinitions = append(validDefinitions, def)
		}
	}
	return validDefinitions
}
