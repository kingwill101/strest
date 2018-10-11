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

//LaunchRequest launch individual request
func LaunchRequest(k string, v strest.Request, validator *validators.Validator) {

	requestStartTime := time.Now()

	logger := log.WithFields(logrus.Fields{
		"request": strings.TrimRight(k, "\n"),
		"start":   requestStartTime,
	})
	// log.Info("Running [%s]\n", k)
	fmt.Printf("Running [%s]\n", k)
	var repeat int
	if v.Repeat == 0 {
		repeat = 1
	} else {
		repeat = v.Repeat
	}

	var loopWg sync.WaitGroup

	for i := 0; i < repeat; i++ {

		loopWg.Add(1)
		go func() {
			if v.Delay > 0 {
				strest.GetLogger().Info("Delaying", k, "for", v.Delay)
				time.Sleep(time.Duration(v.Delay) * time.Millisecond)
			}
			r, err := strest.SendRequest(v)

			if err != nil {
				//TODO handle fail on error
				// validator.Validate(&v, r)
				logger.WithFields(logrus.Fields{"End": time.Now()}).Error(fmt.Sprintf("Error running %s", err.Error()))
			} else {
				validator.Validate(&v, r)

			}
			loopWg.Done()

		}()

	}
	loopWg.Wait()
	fmt.Println("finished in ", time.Since(requestStartTime))

}
