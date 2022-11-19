package strest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

const (
	MethodPOST   string = "POST"
	MethodPUT    string = "PUT"
	MethodGET    string = "GET"
	MethodPATCH  string = "PATCH"
	MethodOPTION string = "OPTION"
	MethodHEAD   string = "HEAD"
)

// GetLogger project logger
func GetLogger() *logrus.Logger {

	return &logrus.Logger{
		Out:       os.Stderr,
		Formatter: new(logrus.TextFormatter),
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.InfoLevel,
	}

}

// LoadYamlData load from yaml file
func LoadYamlData(filename string) (Payload, error) {
	payload := Payload{}

	request, err := os.ReadFile(filename)

	if err != nil {
		GetLogger().Fatalf("Failed to open %s - %s", filename, err.Error())
		return Payload{}, err
	}

	err = yaml.Unmarshal(request, &payload)
	if err != nil {
		log.Printf("Error marshaling : %v", err)
		return payload, err
	}

	return payload, nil
}

// SendRequest make http request
func SendRequest(r Request) (*http.Response, error) {

	var method string
	var req *http.Request
	var err error

	switch ParseField(strings.ToLower(r.Method)) {
	case "get":
		method = MethodGET
		break
	case "post":
		method = MethodPOST
		break
	case "patch":
		method = MethodPATCH
		break
	case "put":
		method = MethodPUT
		break
	case "option":
		method = MethodOPTION
		break
	case "head":
		method = MethodHEAD
		break
	default:
		log.Fatal("Only Get and Post supported")
	}

	//timeout
	timeout := time.Duration(r.Timeout) * time.Millisecond

	client := http.Client{
		Timeout: timeout,
	}

	//build url
	urlBuilder, err := url.Parse(ParseField(r.URL))

	if err != nil {
		return nil, err
	}

	//add url parameters
	q := urlBuilder.Query()
	if len(r.Data.Params) > 0 {
		for k, v := range r.Data.Params {
			q.Add(k, ParseField(v.(string)))
		}
		urlBuilder.RawQuery = q.Encode()
	}

	if method == "POST" {
		// return c.Post(url, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
		if len(r.Data.Form) > 0 {
			GetLogger().Info("Form request")

			values := url.Values{}
			for k, v := range r.Data.Form {
				fmt.Println(k, " ", v)
				values[k] = []string{ParseField(v)}
			}

			req, err = post(urlBuilder.String(), strings.NewReader(values.Encode()))

			if err != nil {
				return &http.Response{}, err
			}

			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		} else if len(r.Data.JSON) > 0 {

			marshal, jsonError := json.Marshal(r.Data.JSON)

			if jsonError != nil {
				return &http.Response{}, jsonError
			}

			req, err = post(
				urlBuilder.String(), strings.NewReader(string(marshal)))

			if err != nil {
				return &http.Response{}, err
			}
			req.Header.Set("Content-Type", "application/marshal")

		} else {

			req, err = post(urlBuilder.String(), bytes.NewReader([]byte(ParseField(r.Data.Raw))))

			if err != nil {
				return &http.Response{}, err
			}

		}
	} else if method == "GET" {
		req, err = http.NewRequest(method, urlBuilder.String(), nil)
		if err != nil {
			return &http.Response{}, err
		}
	} else {

		req, err = http.NewRequest(method, urlBuilder.String(), bytes.NewReader([]byte(r.Data.Raw)))
		if err != nil {
			return &http.Response{}, err
		}
	}

	//add headers
	if len(r.Data.Headers) > 0 {
		for k, v := range r.Data.Headers {
			req.Header.Set(k, ParseField(v))
		}

	}

	resp, err := client.Do(req)
	return resp, err

}

func post(url string, body io.Reader) (req *http.Request, err error) {
	req, err = http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	return
}

// ReadBody process response body
func ReadBody(r *http.Response) (string, error) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(r.Body)
	var bodyBytes []byte
	var err error

	bodyBytes, err = io.ReadAll(r.Body)

	if err != nil {
		log.Println("", err)
		return "", fmt.Errorf("error parsing request body - %s", err.Error())
	}

	// Restore the io.ReadCloser to its original state
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	return string(bodyBytes), nil

}

// RequestKeyExist search through list of requests provided
func RequestKeyExist(key string, p Payload) bool {
	if _, ok := p.Request[key]; !ok {
		return false
	}
	return true
}
