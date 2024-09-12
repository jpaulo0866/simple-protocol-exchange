package main

import (
	"os"

	"gopkg.in/yaml.v2"
)

// Estrutura para o arquivo YAML
type RouteConfig struct {
	Routes []Route `yaml:"routes"`
}

type Route struct {
	Protocol  string      `yaml:"protocol"`
	Port      int         `yaml:"port"`
	Name      string      `yaml:"name"`
	Entry     EntryConfig `yaml:"entry"`
	Transform Transform   `yaml:"transform"`
	Output    Output      `yaml:"output"`
}

type EntryConfig struct {
	BasePath    string `yaml:"basePath"`
	ContentType string `yaml:"content_type"`
	Compressed  bool   `yaml:"compressed"`
}

type Transform struct {
	Remap        []FieldMap        `yaml:"remap"`
	StaticFields map[string]string `yaml:"static_fields"`
	RemoveFields []string          `yaml:"remove_fields"`
}

type FieldMap struct {
	Source         string `yaml:"source"`
	Target         string `yaml:"target"`
	PreserveSource bool   `yaml:"preserve_source"`
}

type Output struct {
	Host        string            `yaml:"host"`
	Port        int               `yaml:"port"`
	Protocol    string            `yaml:"protocol"`
	Path        string            `yaml:"path"`
	Timeout     int               `yaml:"timeout"`
	Headers     map[string]string `yaml:"headers"`
	FilePattern string            `yaml:"file_pattern"`
}

func parseConfig(filePath string) (*RouteConfig, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var config RouteConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
