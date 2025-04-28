package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type IbnSearchFilter struct {
	Name string `json:"name"`
}

type IbnSearch struct {
	PageNumber int             `json:"page-number"`
	PageSize   int             `json:"page-size"`
	Filter     IbnSearchFilter `json:"filter,omitempty"`
}

type IbnSearchPayload struct {
	Input IbnSearch `json:"ibn-administration:input"`
}

type IntentType struct {
	Name    string `json:"name"`
	Version int    `json:"version"`
}

type IbnSearchResponse struct {
	Output struct {
		PageSize   int          `json:"page-size"`
		TotalCount int          `json:"total-count"`
		IntentType []IntentType `json:"intent-type"`
	} `json:"ibn-administration:output"`
}

type IntentTypeYangModule struct {
	Name        string `json:"name"`
	YangContent string `json:"yang-content"`
}

type IntentTypeDefinition struct {
	IntentType struct {
		Module []IntentTypeYangModule `json:"module"`
	} `json:"ibn-administration:intent-type"`
}

// Get available NSP intent types with pagination
func (s *srv) intentTypeSearch(pageNumber, pageSize int, nameFilter string) ([]string, int, error) {
	payload := IbnSearchPayload{
		Input: IbnSearch{
			PageNumber: pageNumber,
			PageSize:   pageSize,
		},
	}

	if nameFilter != "" {
		payload.Input.Filter.Name = nameFilter
	}

	reqBody, err := json.Marshal(payload)
	if err != nil {
		return nil, 0, fmt.Errorf("[Error] generating intent type search payload: %s", err)
	}

	url := fmt.Sprintf("https://%s/restconf/operations/ibn-administration:search-intent-types", s.nsp.Ip)
	resp, err := s.makeHTTPRequest("POST", url, bytes.NewReader(reqBody), nil)
	if err != nil {
		return nil, 0, fmt.Errorf("[Error] fetching intent types: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, 0, fmt.Errorf("[Error] fetching intent types, status: %d", resp.StatusCode)
	}

	var ibnSearchResponse IbnSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&ibnSearchResponse); err != nil {
		return nil, 0, fmt.Errorf("[Error] decoding intent type search response: %s", err)
	}

	var intentTypes []string
	for _, intentType := range ibnSearchResponse.Output.IntentType {
		intentTypes = append(intentTypes, fmt.Sprintf("%s_v%d", intentType.Name, intentType.Version))
	}

	// Handle pagination
	/*if ibnSearchResponse.Output.PageSize == pageSize && ibnSearchResponse.Output.TotalCount > pageSize {
		nextPage, err := s.intentTypeSearch(pageNumber+1, pageSize)
		if err != nil {
			return nil, 0, fmt.Errorf("[Error] fetching paginated intent types: %s", err)
		}
		intentTypes = append(intentTypes, nextPage...)
	}*/

	return intentTypes, ibnSearchResponse.Output.TotalCount, nil
}

// Get YANG modules for a specific NSP intent type
func (s *srv) intentTypeYangModules(intentType string) ([]IntentTypeYangModule, error) {
	lastInd := strings.LastIndex(intentType, "_")
	if lastInd == -1 {
		return nil, fmt.Errorf("[Error] invalid intent type format: %s", intentType)
	}

	name := intentType[:lastInd]
	version := intentType[lastInd+2:]

	url := fmt.Sprintf(
		"https://%s/restconf/data/ibn-administration:ibn-administration/intent-type-catalog/intent-type=%s,%s",
		s.nsp.Ip, url.QueryEscape(name), version,
	)

	resp, err := s.makeHTTPRequest("GET", url, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("[Error] fetching YANG modules for intent type: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("[Error] fetching YANG modules, status: %d", resp.StatusCode)
	}

	var yangDef IntentTypeDefinition
	if err := json.NewDecoder(resp.Body).Decode(&yangDef); err != nil {
		return nil, fmt.Errorf("[Error] decoding YANG modules response: %s", err)
	}

	return yangDef.IntentType.Module, nil
}

// NSP INTENT EXPLORER
func (s *srv) intentExplorer(w http.ResponseWriter, r *http.Request) {
	type PostData struct {
		Url       string `json:"url"`
		Target    string `json:"target"`
		IntentKey string `json:"intent-key"`
	}

	type NspError struct {
		Error any `json:"error"`
	}

	type NspFindErrorRepeat struct {
		RestconfError NspError `json:"ietf-restconf:errors"`
	}

	type NspFindError struct {
		RestconfError NspFindErrorRepeat `json:"ietf-restconf:errors"`
	}

	var pd PostData
	if err := json.NewDecoder(r.Body).Decode(&pd); err != nil {
		s.raiseError("decoding NSP find request failed", err, w)
		return
	}

	url := fmt.Sprintf("https://%s/mdt/rest/restconf/data/ibn:ibn/intent=%s,%s/intent-specific-data%s", s.nsp.Ip, pd.Target, pd.IntentKey, pd.Url)
	headers := map[string]string{
		"Content-Type":  "application/yang-data+json",
		"Accept":        "application/yang-data+json",
		"Authorization": "Bearer " + s.nsp.token.AccessToken,
	}
	resp, err := s.makeHTTPRequest("GET", url, nil, headers)
	if err != nil {
		s.raiseError("error fetching NSP intent explorer request", err, w)
	}

	var output any
	if resp.StatusCode == 200 {
		if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
			s.raiseError("error decoding NSP intent explorer success response", err, w)
		}
	} else {
		var errorResponse NspFindError
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			s.raiseError("error decoding NSP intent explorer error response", err, w)
		}
		output = errorResponse.RestconfError.RestconfError.Error
	}

	response, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		s.raiseError("error creating JSON", err, w)
		return
	}

	writeJsonResponse(w, response)
}
