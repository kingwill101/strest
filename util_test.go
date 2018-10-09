package strest

import (
	"testing"
)

func getYaml() (Payload, error) {
	return LoadYamlData("test/test.yaml")
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
func TestJsonMarshal(t *testing.T) {

}

func TestRequestKeyExist(t *testing.T) {
	y, err := getYaml()

	if err != nil {
		t.Error(err)
	}

	type args struct {
		key string
		p   Payload
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "notEmptyValidation",
			args: args{
				key: "notEmptyValidation",
				p:   y,
			},
			want: true,
		}, {
			name: "helloWorld",
			args: args{
				key: "helloWorld",
				p:   y,
			},
			want: true,
		},
		{
			name: "fake",
			args: args{
				key: "fake",
				p:   y,
			},
			want: true,
		},
		{
			name: "nonexistingkey",
			args: args{
				key: "Blablaj",
				p:   y,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RequestKeyExist(tt.args.key, tt.args.p); got != tt.want {
				t.Errorf("RequestKeyExist() = %v, want %v", got, tt.want)
			}
		})
	}
}
