package validators

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gitlab.com/kingwill101/strest"
	"net/http"
	"strings"
)

// BodyValidator validate status
type BodyValidator struct {
}

// Validate body
func (s BodyValidator) Validate(sreq *strest.Request, res *http.Response, logger *logrus.Entry) bool {
	entry := logger.WithField("validator", "body")
	body, err := strest.ReadBody(res)
	if err != nil {
		entry.Fatal(err.Error())
		return false
	}
	body = strings.TrimRight(body, "\n")

	for k, v := range sreq.Validation {
		entryLog := entry.WithFields(logrus.Fields{
			"validation_tag": k,
		})

		if v.Body == body {
			entryLog.Infof("checks out -  expected %s  got %s", v.Body, body)
		} else {
			failureMsg := fmt.Sprintf("body validation failed expected %s  got %s", v.Body, body)
			if sreq.FailOnError {
				entryLog.Fatal(failureMsg)
			}

			entryLog.Error(failureMsg)
		}
	}
	return true
}
