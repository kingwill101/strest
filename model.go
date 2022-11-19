package strest

import (
	"encoding/json"
	"os"
)

// Payload entire yaml configuration
type Payload struct {
	Version string              `yaml:"version"`
	Request map[string]*Request `yaml:"requests"`
	Async   bool                `yaml:"async"`
}

// RequestData Request data parameters
type RequestData struct {
	Params  map[string]interface{} `yaml:"params"`
	Headers map[string]string      `yaml:"headers"`
	Form    map[string]string      `yaml:"form"`
	JSON    map[string]interface{} `yaml:"json"`
	Raw     string                 `yaml:"raw"`
}

// Request Information about the current request
type Request struct {
	FailOnError bool                   `yaml:"failOnError"`
	URL         string                 `yaml:"url"`
	Timeout     int                    `yaml:"timeout"`
	Method      string                 `yaml:"method"`
	Data        *RequestData           `yaml:"data"`
	Log         bool                   `yaml:"log"`
	Validation  map[string]*Validation `yaml:"validation"`
	Repeat      int                    `yaml:"repeat"`
	Delay       int                    `yaml:"delay"`
	Async       bool                   `yaml:"async"`
	Env         map[string]interface{} `yaml:"environment"`
	DependsOn   []string               `yaml:"dependsOn"`
}

// Validation types of validation available on each request
type Validation struct {
	Body       string `yaml:"body"`
	StatusCode int    `yaml:"code"`
}
