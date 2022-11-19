package main

import (
	"github.com/spf13/cobra"
	"gitlab.com/kingwill101/strest"
	"gitlab.com/kingwill101/strest/runner"
	"gitlab.com/kingwill101/strest/validators"
	"time"
)

func main() {

	var script string

	var rootCmd = &cobra.Command{
		Use:   "strest",
		Short: "run http tests",
		Long:  `Run http tests described with a predefined yaml template`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {

			log := strest.GetLogger()

			validator := validators.NewValidator()
			validator.Register(&validators.StatusCodeValidator{})
			validator.Register(&validators.BodyValidator{})
			validator.Register(&validators.LogPrint{})
			validator.Register(&validators.LogIncoming{})

			testTime := time.Now()

			p, err := strest.LoadYamlData(script)
			if err != nil {
				log.Fatal(err)
			}

			p.Load()

			runner.RunTest(validator, p)

			log.Println("Completed in - ", time.Since(testTime))
		},
	}
	rootCmd.Flags().StringVar(&script, "file", "../../test/req.yaml", "Path to yaml")

	err := rootCmd.Execute()
	if err != nil {
		return
	}
}
