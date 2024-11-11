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

// for intentTypeSearch
type IbnInput struct {
	PageNumber int `json:"page-number"`
	PageSize   int `json:"page-size"`
}

type IntentTypeSearchPayload struct {
	Input IbnInput `json:"ibn-administration:input"`
}

type IntentType struct {
	Name    string `json:"name"`
	Version int    `json:"version"`
}

type IntentTypeList struct {
	Output struct {
		PageSize   int          `json:"page-size"`
		TotalCount int          `json:"total-count"`
		IntentType []IntentType `json:"intent-type"`
	} `json:"ibn-administration:output"`
}

// for intentTypeYangModules
type IntentTypeYangModule struct {
	Name        string `json:"name"`
	YangContent string `json:"yang-content"`
}

type IntentTypeDefinition struct {
	IntentType struct {
		Module []IntentTypeYangModule `json:"module"`
	} `json:"ibn-administration:intent-type"`
}

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
	s.logger.Printf("[Success] NSP Access reset")
}

// getToken
func (s *srv) getToken() error {
	payload := map[string]string{
		"grant_type": "client_credentials",
	}

	auth := s.nsp.User + ":" + s.nsp.Pass
	authEncoded := base64.StdEncoding.EncodeToString([]byte(auth))

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("[Error] parsing NSP access payload: %v", err)
	}

	url := fmt.Sprintf("https://%s/rest-gateway/rest/api/v1/auth/token", s.nsp.Ip)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("[Error] creating NSP access request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic "+authEncoded)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   10 * time.Second,
	}

	resp, err := client.Do(req)
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

// revokeToken
func (s *srv) revokeToken() error {
	payload := map[string]string{
		"token":           s.nsp.token.AccessToken,
		"token_type_hint": "token",
	}

	encodedData := url.Values{}
	for key, value := range payload {
		encodedData.Set(key, value)
	}

	apiBody := strings.NewReader(encodedData.Encode())

	auth := s.nsp.User + ":" + s.nsp.Pass
	authEncoded := base64.StdEncoding.EncodeToString([]byte(auth))

	url := fmt.Sprintf("https://%s/rest-gateway/rest/api/v1/auth/revocation", s.nsp.Ip)
	req, err := http.NewRequest("POST", url, apiBody)
	if err != nil {
		return fmt.Errorf("[Error] creating NSP revocation request: %v", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic "+authEncoded)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("[Error] creating NSP revocation client: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("[Error] revoking NSP access: %v", resp.Body)
	}

	defer s.nspReset()
	return nil
}

// intentTypeSearch
func (s *srv) intentTypeSearch(pageNumber int, pageSize int) ([]string, error) {
	var f []string

	payload := IntentTypeSearchPayload{
		Input: IbnInput{
			PageNumber: pageNumber,
			PageSize:   pageSize,
		},
	}

	buf, err := json.Marshal(payload)
	if err != nil {
		return f, fmt.Errorf("unable to generate payload: %s", err)
	}

	url := fmt.Sprintf("https://%s/restconf/operations/ibn-administration:search-intent-types", s.nsp.Ip)
	req, _ := http.NewRequest("POST", url, bytes.NewReader(buf))

	req.Header.Set("Accept", "application/yang-data+json")
	req.Header.Set("Content-Type", "application/yang-data+json")
	req.Header.Set("Authorization", "Bearer "+s.nsp.token.AccessToken)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   5 * time.Second,
	}

	res, err := client.Do(req)
	if err != nil {
		return f, fmt.Errorf("[Error] creating NSP Intent Type Search client: %s", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return f, fmt.Errorf("[Error] fetching NSP Intent Type Search: %s", err)
	}

	var nspList IntentTypeList
	err = json.NewDecoder(res.Body).Decode(&nspList)
	if err != nil {
		return f, fmt.Errorf("[Error] decoding NSP Intent Type Search response: %s", err)
	}

	for _, intentType := range nspList.Output.IntentType {
		f = append(f, fmt.Sprintf("%s_%d", intentType.Name, intentType.Version))
	}

	if nspList.Output.PageSize == pageSize && nspList.Output.TotalCount > pageSize {
		intentTypeList, err := s.intentTypeSearch(pageNumber+1, pageSize)
		if err != nil {
			return f, fmt.Errorf("[Error] fetching NSP Intent Type Search paginated list: %s", err)
		}
		f = append(f, intentTypeList...)
	}

	return f, nil
}

// intentTypeYangModules
func (s *srv) intentTypeYangModules(intentType string) ([]IntentTypeYangModule, error) {
	var f []IntentTypeYangModule
	fmt.Println(intentType)

	lastInd := strings.LastIndex(intentType, "_")
	name := intentType[:lastInd]
	version := intentType[lastInd+1:]

	url := fmt.Sprintf("https://%s/restconf/data/ibn-administration:ibn-administration/intent-type-catalog/intent-type=%s,%s", s.nsp.Ip, url.QueryEscape(name), version)
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Set("Accept", "application/yang-data+json")
	req.Header.Set("Content-Type", "application/yang-data+json")
	req.Header.Set("Authorization", "Bearer "+s.nsp.token.AccessToken)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   10 * time.Second,
	}

	res, err := client.Do(req)
	if err != nil {
		return f, fmt.Errorf("[Error] creating NSP Intent Type Yang Modules client: %s", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return f, fmt.Errorf("[Error] fetching NSP Intent Type Yang Modules: %s", err)
	}

	var yangDef IntentTypeDefinition
	err = json.NewDecoder(res.Body).Decode(&yangDef)
	if err != nil {
		return f, fmt.Errorf("[Error] decoding NSP Intent Type Yang Modules response: %s", err)
	}

	return yangDef.IntentType.Module, nil
}
