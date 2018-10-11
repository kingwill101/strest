package runner

import (
	"testing"

	"gitlab.com/kingwill101/strest"
	"gitlab.com/kingwill101/strest/validators"
)

func getYaml() (strest.Payload, error) {
	return strest.LoadYamlData("../test/test.yaml")
}

func getValidators() *validators.Validator {
	validator := validators.NewValidator()
	// validator.Register(&validators.StatusCodeValidator{})
	// validator.Register(&validators.StatusValidator{})
	// validator.Register(&validators.LogPrint{})
	validator.Register(&validators.LogIncoming{})

	return validator
}
func TestRunTest(t *testing.T) {
	v := getValidators()
	p, err := getYaml()

	if err != nil {
		t.Error(err)
	}

	type args struct {
		validator *validators.Validator
		p         strest.Payload
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test runner",
			args: args{
				validator: v,
				p:         p,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RunTest(tt.args.validator, tt.args.p)
		})
	}
}
