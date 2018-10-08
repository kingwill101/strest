package strest

import "testing"

func getYaml() (Payload, error) {
	return LoadYamlData("test/req.yaml")
}
func TestLoadYaml(t *testing.T) {
	_, err := getYaml()
	if err != nil {
		t.Error("Unable to load yaml config file - ", err)
	}
}

func TestEmptyValidation(t *testing.T) {
	y, err := getYaml()

	if err != nil {
		t.Error(err)
	}

	if _, ok := y.Request["emptyValidation"]; !ok {
		t.Error("key emptyValidation is missing")
	}

	if (Validation{}) != y.Request["emptyValidation"].Validation {
		t.Error("Validation should be empty")
	}

}
func TestNotEmptyValidation(t *testing.T) {
	y, err := getYaml()

	if err != nil {
		t.Error(err)
	}
	if _, ok := y.Request["notEmptyValidation"]; !ok {
		t.Error("key notEmptyValidation is missing")
	}

	if (Validation{}) == y.Request["notEmptyValidation"].Validation {
		t.Error("Validation should not be empty")
	}

}
