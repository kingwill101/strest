package strest

import (
	"bytes"
	"os"
	"text/template"
)

//ParseField parse fields that need access to template functions
func ParseField(in string) string {
	tmp, err := getTemplate().Parse(in)
	if err != nil {
		log.Fatal(err)
	}
	buf := bytes.Buffer{}

	err = tmp.Execute(&buf, nil)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println("PARSE  IN -->", in)
	// fmt.Println("PARSE OUT -->", buf.String())

	return buf.String()
}

func getTemplate() *template.Template {
	return template.New("strest").Funcs(tmpFuncs())
}

func tmpFuncs() template.FuncMap {
	return template.FuncMap{
		"ENV": getEnv,
	}
}

func getEnv(in string) string {
	log.Println("Searching ", in, " --> ", os.Getenv(in))
	return os.Getenv(in)
}
