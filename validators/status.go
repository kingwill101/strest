package validators

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	strest "gitlab.com/kingwill101/strest"
)

//StatusValidator validate status
type StatusValidator struct {
}

func (s StatusValidator) Validate(sreq *strest.Request, res *http.Response) bool {
	if sreq.Validation.StatusCode == 0 {
		return false
	}

	fmt.Println("[status code]")

	if sreq.Validation.StatusCode != res.StatusCode {

		fmt.Printf("\tGot - %d\n\tExpected - %d \n", res.StatusCode, sreq.Validation.StatusCode)
		return false

	} else {
		fmt.Println("\tSuccess [status] checks out")
	}
	return true
}

//BodyValidator validate status
type BodyValidator struct {
}

func (s BodyValidator) Validate(sreq *strest.Request, res *http.Response) bool {
	fmt.Println("[body]")
	body, err := strest.ReadBody(res)
	if err != nil {
		log.Fatal(err.Error())
		return false
	}
	body = strings.TrimRight(body, "\n")
	if sreq.Validation.Body != body {

		fmt.Printf("\tGot - %s\n\tExpected - %s \n", body, sreq.Validation.Body)
		return false

	} else {
		fmt.Println("\tSuccess [body] checks out")
	}
	return true
}

//StatusCodeValidator validate status
type StatusCodeValidator struct {
}

func (s StatusCodeValidator) Validate(sreq *strest.Request, res *http.Response) bool {
	fmt.Println("[body]")
	body, err := strest.ReadBody(res)
	if err != nil {
		log.Fatal(err.Error())
		return false
	}
	body = strings.TrimRight(body, "\n")
	if sreq.Validation.Body != body {

		fmt.Printf("\tGot - %s\n\tExpected - %s \n", body, sreq.Validation.Body)
		return false

	} else {
		fmt.Println("\tSuccess [body] checks out")
	}
	return true
}

func BodyValidateFunc(sreq *strest.Request, res *http.Response) bool {
	fmt.Println("[body]")
	body, err := strest.ReadBody(res)
	if err != nil {
		log.Fatal(err.Error())
		return false
	}
	body = strings.TrimRight(body, "\n")
	if sreq.Validation.Body != body {

		fmt.Printf("\tGot - %s\n\tExpected - %s \n", body, sreq.Validation.Body)
		return false

	} else {
		fmt.Println("\tSuccess [body] checks out")
	}
	return true
}
