package runner

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"gitlab.com/kingwill101/strest"
	"gitlab.com/kingwill101/strest/validators"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func init() {
	log = strest.GetLogger()
}

//RunTest execute all test
func RunTest(validator *validators.Validator, p strest.Payload) {

	var wg sync.WaitGroup

	for k, v := range p.Request {

		wg.Add(1)

		go func(k string, v strest.Request, validator *validators.Validator) {
			LaunchRequest(k, v, validator)
			wg.Done()
		}(k, v, validator)
	}

	wg.Wait()
}

// LaunchRequest launch individual request
func LaunchRequest(k string, request strest.Request, validator *validators.Validator) {

	requestStartTime := time.Now()

	logger := strest.GetLogger().WithFields(logrus.Fields{
		"request":     strings.TrimRight(k, "\n"),
		"using async": request.Async,
		"launched":    requestStartTime,
	})

	logger.Infof("Running [%s]", k)
	var repeat int
	if request.Repeat == 0 {
		repeat = 1
	} else {
		repeat = request.Repeat
	}

	var loopWg sync.WaitGroup

	for i := 0; i < repeat; i++ {

		requestSender := func() {

			logger.Debug("request launched ")
			if request.Delay > 0 {
				logger.Debug("Delaying %s for %d", k, request.Delay)
				time.Sleep(time.Duration(request.Delay) * time.Millisecond)
			}
			r, err := strest.SendRequest(request)

			if err != nil {
				//TODO handle fail on error
				// validator.Validate(&request, r)
				logger.WithFields(logrus.Fields{"End": time.Now()}).Error(fmt.Sprintf("Error running %s", err.Error()))
			} else {
				logger.Debug("About to validate")
				validator.Validate(&request, r, logger)
			}

			if request.Async {
				loopWg.Done()
			}
		}

		if request.Async {
			loopWg.Add(1)
			go requestSender()

		} else {
			logger.Debug("Launching ")
			requestSender()
		}
	}

	if request.Async {
		loopWg.Wait()
	}

	fmt.Println("finished in ", time.Since(requestStartTime))

}
