package main

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// REST NSP ACCESS
func (s *srv) nspReset() {
	s.Lock()
	s.nsp = NspAccess{
		Ip:   "",
		User: "",
		Pass: "",
		token: TokenDetail{
			AccessToken:  "",
			RefreshToken: "",
			TokenType:    "",
			ExpiresIn:    0,
		},
	}
	s.Unlock()
	s.logger.Printf("[Info] NSP Access reset")
}

// GET NSP TOKEN
func (s *srv) getToken() error {
	payload := map[string]string{
		"grant_type": "client_credentials",
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("[Error] parsing NSP access payload: %v", err)
	}

	auth := s.nsp.User + ":" + s.nsp.Pass
	authEncoded := base64.StdEncoding.EncodeToString([]byte(auth))

	url := fmt.Sprintf("https://%s/rest-gateway/rest/api/v1/auth/token", s.nsp.Ip)
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Basic " + authEncoded,
	}

	resp, err := s.makeHTTPRequest("POST", url, bytes.NewBuffer(jsonData), headers)
	if err != nil {
		return fmt.Errorf("[Error] creating NSP access client: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("[Error] NSP access request failed with status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("[Error] accessing NSP access response: %v", err)
	}

	err = json.Unmarshal(body, &s.nsp.token)
	if err != nil {
		return fmt.Errorf("[Error] parsing NSP access response: %v", err)
	}

	return nil
}

// REVOKE NSP TOKEN
func (s *srv) revokeToken() error {
	payload := url.Values{
		"token":           {s.nsp.token.AccessToken},
		"token_type_hint": {"token"},
	}
	apiBody := strings.NewReader(payload.Encode())

	auth := s.nsp.User + ":" + s.nsp.Pass
	authEncoded := base64.StdEncoding.EncodeToString([]byte(auth))

	url := fmt.Sprintf("https://%s/rest-gateway/rest/api/v1/auth/revocation", s.nsp.Ip)
	headers := map[string]string{
		"Content-Type":  "application/x-www-form-urlencoded",
		"Authorization": "Basic " + authEncoded,
	}

	resp, err := s.makeHTTPRequest("POST", url, apiBody, headers)
	if err != nil {
		return fmt.Errorf("[Error] creating NSP revocation client: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("[Error] revoking NSP access with status: %s", resp.Status)
	}

	s.nspReset()
	return nil
}

// makeHTTPRequest creates and executes an HTTP request
func (s *srv) makeHTTPRequest(method, url string, body io.Reader, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("[Error] creating HTTP request: %v", err)
	}

	if headers != nil {
		for key, value := range headers {
			req.Header.Set(key, value)
		}
	} else {
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+s.nsp.token.AccessToken)
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: 7 * time.Second,
	}

	return client.Do(req)
}

// GET NSP STATUS
func (s *srv) getNspStatus() error {
	url := fmt.Sprintf("https://%s/nsp-role-manager/rest/api/v1/server/status", s.nsp.Ip)

	resp, err := s.makeHTTPRequest("GET", url, nil, nil)
	if err != nil {
		return fmt.Errorf("[Error] triggering NSP status: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("[Error] fetching NSP health: %d", resp.StatusCode)
	}

	return nil
}

// NSP INVENTORY FIND
func (s *srv) nspFind(w http.ResponseWriter, r *http.Request) {
	type NspFindInput struct {
		XpathFilter string `json:"xpath-filter"`
		IncludeMeta bool   `json:"include-meta"`
		Fields      string `json:"fields,omitempty"`
		Depth       int    `json:"depth"`
		Limit       int    `json:"limit"`
		Offset      int    `json:"offset"`
	}

	type RequestPayload struct {
		Kind string       `json:"kind"`
		Nsp  NspFindInput `json:"nsp"`
	}

	type NspFindOutputData struct {
		StartIndex int `json:"start-index"`
		EndIndex   int `json:"end-index"`
		TotalCount int `json:"total-count"`
		Data       any `json:"data"`
	}

	type NspFindOutput struct {
		Output NspFindOutputData `json:"nsp-inventory:output,omitempty"`
	}

	type NspError struct {
		Error any `json:"error"`
	}

	type NspFindError struct {
		RestconfError NspError `json:"ietf-restconf:errors"`
	}

	type NspFindPayload struct {
		Input NspFindInput `json:"input"`
	}

	var requestPayload RequestPayload
	if err := json.NewDecoder(r.Body).Decode(&requestPayload); err != nil {
		s.raiseError("decoding NSP find request failed", err, w)
		return
	}

	payload := NspFindPayload{
		Input: requestPayload.Nsp,
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

	var output any
	if resp.StatusCode == 200 {
		var successResponse NspFindOutput
		if err := json.NewDecoder(resp.Body).Decode(&successResponse); err != nil {
			s.raiseError("error decoding NSP find success response", err, w)
		}
		output = successResponse.Output
	} else {
		var errorResponse NspFindError
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			s.raiseError("error decoding NSP find error response", err, w)
		}
		output = errorResponse.RestconfError
	}

	response, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		s.raiseError("error creating JSON", err, w)
		return
	}

	writeJsonResponse(w, response)
}
