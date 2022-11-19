package runner

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"gitlab.com/kingwill101/strest"
	"gitlab.com/kingwill101/strest/validators"
	"golang.org/x/exp/slices"
)

// RunTest execute all test
func RunTest(validator *validators.Validator, p strest.Payload) {

	strest.GetLogger().Debug("Warming up test runners ")
	strest.GetLogger().Debug("using async - ", p.Async)

	var called []string

	var requestNames []string

	for k := range p.Request {
		requestNames = append(requestNames, k)
	}

	for {
		if len(called) < len(p.Request) {

			for k, v := range p.Request {

				if slices.Contains(called, k) {
					continue
				}

				if len(v.DependsOn) > 0 {
					dependenciesSatisfied := true

					for _, dependency := range v.DependsOn {
						if !slices.Contains(requestNames, dependency) {
							strest.GetLogger().Warnf("Unknown dependency %s", dependency)
						}
						if slices.Contains(requestNames, dependency) && !slices.Contains(called, dependency) {
							dependenciesSatisfied = false
						} else {
							strest.GetLogger().Infof("Dependency already satisfied %s", dependency)
						}
					}

					if dependenciesSatisfied {
						LaunchRequest(k, *v, validator)
						called = append(called, k)
						println("\n\n\n")
					}

				} else {

					if len(v.DependsOn) == 0 {

						println("\n\n\n")
						LaunchRequest(k, *v, validator)
						called = append(called, k)
					}
				}
			}
		} else {
			break
		}
	}

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
