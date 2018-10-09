package strest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var log *logrus.Logger

func init() {
	log = GetLogger()
}

//GetLogger project logger
func GetLogger() *logrus.Logger {
	// logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	// logrus.SetLevel(logrus.WarnLevel)
	return logrus.New()
}

//LoadYamlData load from yaml file
func LoadYamlData(filename string) (Payload, error) {
	payload := Payload{}

	request, err := ioutil.ReadFile(filename)

	if err != nil {
		log.Fatalf("Failed to open %s - %s", filename, err.Error())
		return Payload{}, err
	}

	err = yaml.Unmarshal([]byte(request), &payload)
	if err != nil {
		log.Printf("Error marshaling : %v", err)
		return payload, err
	}

	return payload, nil
}

//SendRequest make http request
func SendRequest(r Request) (*http.Response, error) {

	var method string
	var req *http.Request
	var err error

	switch ParseField(strings.ToLower(r.Method)) {
	case "get":
		method = "GET"
	case "post":
		method = "POST"
	case "put":
		method = "PUT"
	case "option":
		method = "OPTION"
	default:
		log.Fatal("Only Get and Post supported")
	}

	//timeout
	timeout := time.Duration(time.Duration(r.Timeout) * time.Millisecond)

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
			log.Info("Form request")

			vals := url.Values{}
			for k, v := range r.Data.Form {
				fmt.Println(k, " ", v)
				vals[k] = []string{ParseField(v)}
			}

			req, err = post(urlBuilder.String(), strings.NewReader(vals.Encode()))

			if err != nil {
				return &http.Response{}, err
			}

			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		} else if len(r.Data.JSON) > 0 {

			jso, jerr := json.Marshal(r.Data.JSON)

			if jerr != nil {
				return &http.Response{}, jerr
			}

			req, err = post(
				urlBuilder.String(), strings.NewReader(string(jso)))

			if err != nil {
				return &http.Response{}, err
			}
			req.Header.Set("Content-Type", "application/json")

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
	// req.Header.Set("Content-Type", bodyType)
	return
}

//ReadBody process response body
func ReadBody(r *http.Response) (string, error) {
	defer r.Body.Close()
	var bodyBytes []byte
	var err error

	bodyBytes, err = ioutil.ReadAll(r.Body)

	if err != nil {
		log.Println("", err)
		return "", fmt.Errorf("Error parsing request body - %s", err.Error())
	}

	// Restore the io.ReadCloser to its original state
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	return string(bodyBytes), nil

}
