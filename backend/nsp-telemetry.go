package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

type NodeSupportEntry struct {
	Node     string   `json:"node"`
	Releases []string `json:"releases"`
}

func (s *srv) getTelemetryTypes(w http.ResponseWriter, r *http.Request) {
	// INPUT
	url := fmt.Sprintf("https://%s/restconf/data/telemetry-admin:/ageout-policies/ageout-policy", s.nsp.Ip)
	resp, err := s.makeHTTPRequest("GET", url, nil, nil)
	if err != nil {
		s.raiseError("error fetching telemetry ageout policies", err, w)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		s.raiseError("error querying ageout policies", err, w)
	}

	// OUTPUT
	type TelemetryTypes struct {
		Name string `json:"name"`
	}

	type Policies struct {
		AgeoutPolicy []TelemetryTypes `json:"ageout-policy"`
	}

	var successResponse Policies
	if err := json.NewDecoder(resp.Body).Decode(&successResponse); err != nil {
		s.raiseError("error decoding NSP ageout policies response", err, w)
	}

	var tt []string
	for _, data := range successResponse.AgeoutPolicy {
		tt = append(tt, strings.Replace(data.Name, "telemetry:/", "/", -1))
	}

	response, err := json.MarshalIndent(tt, "", "  ")
	if err != nil {
		s.raiseError("error creating JSON", err, w)
		return
	}

	writeJsonResponse(w, response)
}

func (s *srv) getTelemetryTypeDefinition(w http.ResponseWriter, r *http.Request) {
	saveParam := r.URL.Query().Get("save")
	save := saveParam == "true"

	type DefinitionInput struct {
		Name string `json:"name"`
	}

	var defInput DefinitionInput
	if err := json.NewDecoder(r.Body).Decode(&defInput); err != nil {
		s.raiseError("error decoding telemetry type request payload", err, w)
		return
	}
	telemetryType := defInput.Name

	baseurl := fmt.Sprintf("https://%s/SearchApp/rest/api/v1/documents/telemetryData", s.nsp.Ip)
	// remember order of params matters
	urlParams := fmt.Sprintf("?fq=telemetryType:*telemetry:%s*&start=0&fq=record_type:telemetry_stats_info&rows=500", telemetryType)
	resp, err := s.makeHTTPRequest("POST", baseurl+urlParams, bytes.NewBuffer([]byte("{}")), nil)
	if err != nil {
		s.raiseError("error fetching NSP telemetry stats info", err, w)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		s.raiseError("error querying telemetry stats info", err, w)
		return
	}

	// OUTPUT
	type RowData struct {
		TelemetryType  string             `json:"telemetryType"`
		CounterName    string             `json:"counterName"`
		DataType       string             `json:"dataType"`
		DeviceXpath    string             `json:"deviceXpath"`
		NodeSupport    []NodeSupportEntry `json:"nodeSupport,omitempty"`
		NodeSupportRaw string             `json:"childBody,omitempty"`
	}

	type TelemetryData struct {
		RowData RowData `json:"rowData"`
	}

	type ResponseData struct {
		TelemetryData []TelemetryData `json:"TelemetryData"`
	}

	type Response struct {
		Data ResponseData `json:"data"`
	}

	type DefinitionResponse struct {
		Response Response `json:"response"`
	}

	var successResponse DefinitionResponse
	if err := json.NewDecoder(resp.Body).Decode(&successResponse); err != nil {
		s.raiseError("error querying telemetry stats info", err, w)
		return
	}

	var ttd []RowData
	telemetryData := successResponse.Response.Data.TelemetryData
	for _, data := range telemetryData {
		local := data.RowData
		local.NodeSupport = formatNodeSupportEntries(local.NodeSupportRaw)
		local.NodeSupportRaw = ""
		ttd = append(ttd, local)
	}

	response, err := json.MarshalIndent(ttd, "", "  ")
	if err != nil {
		s.raiseError("error creating JSON", err, w)
		return
	}

	if save {
		filename := strings.ReplaceAll(defInput.Name[1:], "/", "_")
		if err := s.saveJsonFile("telemetry-type", filename, response); err != nil {
			s.raiseError("", err, w)
			return
		}
		writeResponse(w, "success", fmt.Sprintf("telemetry-type/%s.json was saved", filename))
	} else {
		writeJsonResponse(w, response)
	}
}

func formatNodeSupportEntries(input string) []NodeSupportEntry {
	re := regexp.MustCompile(`([^\[\]]+)\s*\[([^\]]+)\]`)

	matches := re.FindAllStringSubmatch(input, -1)

	var results []NodeSupportEntry

	for _, match := range matches {
		kind := strings.TrimSpace(match[1])
		releases := strings.Split(match[2], ",")
		for i := range releases {
			releases[i] = strings.TrimSpace(releases[i])
		}
		results = append(results, NodeSupportEntry{
			Node:     kind,
			Releases: releases,
		})
	}

	return results
}
