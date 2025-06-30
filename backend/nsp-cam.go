package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"sigs.k8s.io/yaml"
)

func (s *srv) getTransformers(tt string) []FlatMapping {
	parts := strings.Split(tt, "/")

	last := parts[len(parts)-1]
	ttHypenated := strings.ReplaceAll(last, "_", "-")

	url := fmt.Sprintf("https://%s/cam/rest/api/v2/artifact/full/", s.nsp.Ip)
	resp, err := s.makeHTTPRequest("GET", url, nil, nil)
	if err != nil {
		fmt.Println("error fetching NSP CAM artifacts", err)
	}

	if resp.StatusCode != 200 {
		fmt.Println("error querying NSP CAM artifacts", err)
	}

	// OUTPUT
	type CamArtifacts struct {
		Artifacts struct {
			Response struct {
				Data []struct {
					Status struct {
						Signature string `json:"signature"`
						Status    string `json:"status"`
						Artifact  struct {
							Author          string `json:"author"`
							Name            string `json:"name"`
							ArtifactName    string `json:"artifactName"`
							ArtifactContent []struct {
								Path     string `json:"path"`
								FileName string `json:"fileName"`
							} `json:"artifactContent"`
						} `json:"artifact"`
					} `json:"status"`
				} `json:"data"`
			} `json:"response"`
		} `json:"artifacts"`
	}

	var cam CamArtifacts
	if err := json.NewDecoder(resp.Body).Decode(&cam); err != nil {
		fmt.Println("error querying NSP CAM artifacts", err)
	}

	for _, entry := range cam.Artifacts.Response.Data {
		flag := false
		if entry.Status.Status == "Installed" &&
			entry.Status.Signature == "NOKIA R&D" &&
			entry.Status.Artifact.Author == "NOKIA R&D" &&
			strings.HasSuffix(entry.Status.Artifact.ArtifactName, "transformer-cr") {
			if (strings.Contains(parts[2], "sros") && strings.HasPrefix(entry.Status.Artifact.ArtifactName, "nsp-sros-ct")) ||
				(!strings.Contains(parts[2], "sros") && strings.HasPrefix(entry.Status.Artifact.ArtifactName, "nsp-sros-va")) {
				flag = true
			}
			if flag {
				for _, file := range entry.Status.Artifact.ArtifactContent {
					if strings.Contains(file.FileName, ttHypenated) {
						return s.getTransformerFile(file.Path, file.FileName)
					}
				}
			}
		}
	}

	return nil
}

func (s *srv) getTransformerFile(path string, fileName string) []FlatMapping {
	url := fmt.Sprintf("https://%s/nsp-file-service-app/rest/api/v1/file/returnContent?filePath=%s/%s", s.nsp.Ip, path, fileName)
	resp, err := s.makeHTTPRequest("GET", url, nil, nil)
	if err != nil {
		fmt.Println("error fetching NSP CAM artifact file definition", err)
	}

	if resp.StatusCode != 200 {
		fmt.Println("error querying NSP CAM artifact file definition", err)
	}

	// OUTPUT
	type TransformerFile struct {
		Data string `json:"data"`
	}

	type YamlDefintion struct {
		Metadata struct {
			Name string `json:"name" yaml:"name"`
		} `json:"metadata" yaml:"metadata"`
		Spec struct {
			NSPPath                string `yaml:"nsp-path"`
			DeviceClassDefinitions []struct {
				Group    string `json:"group" yaml:"group"`
				Mappings []struct {
					DevicePath string `json:"device-path" yaml:"device-path"`
					KeyMatcher struct {
						Path string `json:"path" yaml:"path"`
					} `json:"key-matcher" yaml:"key-matcher"`
				} `json:"mappings" yaml:"mappings"`
			} `json:"deviceclassdefinitions" yaml:"deviceclassdefinitions"`
		} `json:"spec" yaml:"spec"`
	}

	var file TransformerFile
	if err := json.NewDecoder(resp.Body).Decode(&file); err != nil {
		fmt.Println("error querying NSP CAM artifact file definition", err)
	}

	var yamlDef YamlDefintion
	if err := yaml.Unmarshal([]byte(file.Data), &yamlDef); err != nil {
		log.Fatal("YAML error:", err)
	}

	nspPath := yamlDef.Spec.NSPPath
	var flat []FlatMapping

	for _, dc := range yamlDef.Spec.DeviceClassDefinitions {
		for _, m := range dc.Mappings {
			if strings.Contains(m.KeyMatcher.Path, "[") {
				flat = append(flat, FlatMapping{
					NSPPath:    nspPath,
					DevicePath: m.DevicePath,
					Path:       m.KeyMatcher.Path,
				})
			}
		}
	}

	return flat
}
