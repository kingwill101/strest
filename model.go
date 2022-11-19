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

func marshal[K comparable, V any](m map[K]V) (map[K]V, error) {

	stringify, err := json.Marshal(m)

	if err != nil {
		return nil, err
	}

	stringify = []byte(ParseField(string(stringify)))
	var result map[K]V
	if err := json.Unmarshal(stringify, &result); err != nil {
		return nil, err
	}

	return result, nil

}

func (p *Payload) Load() {

	for s, v := range p.Request {
		GetLogger().Println("request", s)

		for env, value := range v.Env {
			err := os.Setenv(env, value.(string))
			if err != nil {
				continue
			}
		}

		p.Request[s].URL = ParseField(v.URL)
		p.Request[s].Method = ParseField(v.Method)
		p.Request[s].Data.Raw = ParseField(v.Data.Raw)

		if len(v.Data.Headers) != 0 {

			headers, err := marshal[string, string](v.Data.Headers)
			if err != nil {
				continue
			}
			p.Request[s].Data.Headers = headers
		}

		if len(v.Data.Form) != 0 {

			form, err := marshal[string, string](v.Data.Form)
			if err != nil {
				continue
			}
			p.Request[s].Data.Headers = form
		}
		if len(v.Data.JSON) != 0 {

			requestJson, err := marshal[string, any](v.Data.JSON)
			if err != nil {
				continue
			}
			p.Request[s].Data.JSON = requestJson
		}

		for t, u := range v.Validation {
			v.Validation[t].Body = ParseField(u.Body)
		}

	}
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
