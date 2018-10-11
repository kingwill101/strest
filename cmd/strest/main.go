package main

import (
	"flag"
	"fmt"
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

	// fmt.Println("ADDRESS", os.GetEnv("ADDRESS"))
	validator := validators.NewValidator()
	// validator.Register(&validators.StatusCodeValidator{})
	// validator.Register(&validators.StatusValidator{})
	// validator.Register(&validators.LogPrint{})
	validator.Register(&validators.LogIncoming{})

	// validator.RegisterFunc(validators.BodyValidateFunc)

	testTime := time.Now()

	p, err := strest.LoadYamlData(*script)
	if err != nil {
		log.Fatal(err)
	}

	runner.RunTest(validator, p)

	fmt.Println("\nCompleted in - ", time.Since(testTime))

	// runServer()
}
