package main

import (
	"flag"
	"time"

	"github.com/sirupsen/logrus"

	"gitlab.com/kingwill101/strest"
	"gitlab.com/kingwill101/strest/runner"
	"gitlab.com/kingwill101/strest/validators"
)

var script = flag.String("script", "../../test/req.yaml", "Yaml file for running validation")
var log *logrus.Logger

func init() {
	log = strest.GetLogger()
}

func main() {

	validator := validators.NewValidator()
	validator.Register(&validators.StatusCodeValidator{})
	validator.Register(&validators.BodyValidator{})
	validator.Register(&validators.LogPrint{})
	validator.Register(&validators.LogIncoming{})

	testTime := time.Now()

	p, err := strest.LoadYamlData(*script)
	if err != nil {
		log.Fatal(err)
	}

	p.Load()

	runner.RunTest(validator, p)

	log.Println("Completed in - ", time.Since(testTime))

	// runServer()
}
