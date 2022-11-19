package runner

import (
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"gitlab.com/kingwill101/strest"
	"gitlab.com/kingwill101/strest/validators"
)

func getYaml() (strest.Payload, error) {
	return strest.LoadYamlData("../test/test.yaml")
}

func getMainEngine() *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/test")
	{
		v1.POST("/user", func(c *gin.Context) {

			c.JSON(200, gin.H{
				"username": "kingwill",
				"password": "some password",
			})
		})
		v1.POST("/helloworld", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"username": "kingwill",
				"password": "some password",
			})
		})

	}
	return r
}

func getValidators() *validators.Validator {
	validator := validators.NewValidator()
	// validator.Register(&validators.StatusCodeValidator{})
	// validator.Register(&validators.StatusValidator{})
	// validator.Register(&validators.LogPrint{})
	validator.Register(&validators.BodyValidator{})
	// validator.Register(&validators.LogIncoming{})

	return validator
}
func TestRunTest(t *testing.T) {
	ts := httptest.NewServer(getMainEngine())
	defer ts.Close()
	_ = os.Setenv("SERVER", ts.URL)

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
