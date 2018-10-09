package strest

import (
	"fmt"
	"os"
	"testing"
)

func TestParseField(t *testing.T) {

	if err := os.Setenv("STREST", "STREST"); err != nil {
		t.Fatal(err)
	}

	type args struct {
		in string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "environ",
			args: args{in: `{{ ENV "HOME" }}`},
			want: "/home/kingwill101",
		},
		{
			name: "environ2",
			args: args{in: `{{ ENV "STREST" }}`},
			want: "STREST",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseField(tt.args.in); got != tt.want {
				t.Errorf("ParseField() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParse(t *testing.T) {
	fmt.Println(ParseField(`{{ ENV "HOME" }}`))
}
