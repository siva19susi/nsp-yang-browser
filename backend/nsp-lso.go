package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (s *srv) getLsoOperations(w http.ResponseWriter, r *http.Request) {
	// INPUT
	type NspFindInput struct {
		XpathFilter string `json:"xpath-filter"`
		IncludeMeta bool   `json:"include-meta"`
		Fields      string `json:"fields,omitempty"`
		Depth       int    `json:"depth"`
	}

	type NspFindPayload struct {
		Input NspFindInput `json:"input"`
	}

	payload := NspFindPayload{
		Input: NspFindInput{
			XpathFilter: "/nsp-lso-manager:lso-manager/operation-types/operation-type",
			IncludeMeta: false,
			Fields:      "name;version",
			Depth:       2,
		},
	}
	reqBody, err := json.Marshal(payload)
	if err != nil {
		s.raiseError("error parsing NSP find payload", err, w)
		return
	}

	url := fmt.Sprintf("https://%s/restconf/operations/nsp-inventory:find", s.nsp.Ip)
	resp, err := s.makeHTTPRequest("POST", url, bytes.NewReader(reqBody), nil)
	if err != nil {
		s.raiseError("error fetching NSP find request", err, w)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		s.raiseError("error querying Operation Types", err, w)
	}

	// OUTPUT
	type OperationTypeList struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	}

	type NspFindOutputData struct {
		StartIndex int                 `json:"start-index"`
		EndIndex   int                 `json:"end-index"`
		TotalCount int                 `json:"total-count"`
		Data       []OperationTypeList `json:"data"`
	}

	type NspFindOutput struct {
		Output NspFindOutputData `json:"nsp-inventory:output,omitempty"`
	}

	var successResponse NspFindOutput
	if err := json.NewDecoder(resp.Body).Decode(&successResponse); err != nil {
		s.raiseError("error decoding NSP find success response", err, w)
	}

	var operationTypes []string
	for _, data := range successResponse.Output.Data {
		operationTypes = append(operationTypes, fmt.Sprintf("%s_v%s", data.Name, data.Version))
	}

	response, err := json.MarshalIndent(operationTypes, "", "  ")
	if err != nil {
		s.raiseError("error creating JSON", err, w)
		return
	}

	writeJsonResponse(w, response)
}

func (s *srv) getLsoOperationModel(operation string) (string, string, error) {
	// COMMON
	fileName := ""
	fileContent := ""

	lastInd := strings.LastIndex(operation, "_")
	if lastInd == -1 {
		return fileName, fileContent, fmt.Errorf("invalid operation type format")
	}

	name := operation[:lastInd]
	version := operation[lastInd+2:]

	// INPUT
	type NspFindInput struct {
		XpathFilter string `json:"xpath-filter"`
		IncludeMeta bool   `json:"include-meta"`
		Fields      string `json:"fields,omitempty"`
		Depth       int    `json:"depth"`
	}

	type NspFindPayload struct {
		Input NspFindInput `json:"input"`
	}

	type OperationTypeList struct {
		OperationModel string `json:"operation-model"`
	}

	payload := NspFindPayload{
		Input: NspFindInput{
			XpathFilter: fmt.Sprintf("/nsp-lso-manager:lso-manager/operation-types/operation-type[name = '%s' and version='%s']", name, version),
			IncludeMeta: false,
			Fields:      "operation-model",
			Depth:       2,
		},
	}
	reqBody, err := json.Marshal(payload)
	if err != nil {
		return fileName, fileContent, err
	}

	url := fmt.Sprintf("https://%s/restconf/operations/nsp-inventory:find", s.nsp.Ip)
	resp, err := s.makeHTTPRequest("POST", url, bytes.NewReader(reqBody), nil)
	if err != nil {
		return fileName, fileContent, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fileName, fileContent, err
	}

	// OUTPUT
	type NspFindOutputData struct {
		StartIndex int                 `json:"start-index"`
		EndIndex   int                 `json:"end-index"`
		TotalCount int                 `json:"total-count"`
		Data       []OperationTypeList `json:"data"`
	}

	type NspFindOutput struct {
		Output NspFindOutputData `json:"nsp-inventory:output,omitempty"`
	}

	var successResponse NspFindOutput
	if err := json.NewDecoder(resp.Body).Decode(&successResponse); err != nil {
		return fileName, fileContent, err
	}

	fileName = name
	fileContent = successResponse.Output.Data[0].OperationModel

	return fileName, fileContent, nil
}
