package validators

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gitlab.com/kingwill101/strest"
	"net/http"
)

// StatusCodeValidator validate status
type StatusCodeValidator struct {
}

// Validate status code
func (s StatusCodeValidator) Validate(sreq *strest.Request, res *http.Response, logger *logrus.Entry) bool {

	entry := logger.WithField("validator", "status-code")
	for k, v := range sreq.Validation {
		entryLog := entry.WithFields(logrus.Fields{
			"validation_tag": k,
		})
		if v.StatusCode == 0 {
			continue
		}
		if v.StatusCode == res.StatusCode {
			entryLog.Infof("checks out -  expected %d  got %d", v.StatusCode, res.StatusCode)
		} else {
			failureMsg := fmt.Sprintf("status validation failed expected %d  got %d", v.StatusCode, res.StatusCode)
			if sreq.FailOnError {
				entryLog.Fatal(failureMsg)
			}

			entryLog.Error(failureMsg)
		}
	}

	return true
}
