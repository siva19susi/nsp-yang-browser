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
