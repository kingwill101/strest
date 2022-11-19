package strest

import (
	"bytes"
	"os"
	"text/template"

	"github.com/icrowley/fake"
)

// ParseField parse fields that need access to template functions
func ParseField(in string) string {
	tmp, err := getTemplate().Parse(in)
	if err != nil {
		GetLogger().Fatal(err)
	}
	buf := bytes.Buffer{}

	err = tmp.Execute(&buf, nil)
	if err != nil {
		GetLogger().Fatal(err)
	}

	return buf.String()
}

func getTemplate() *template.Template {
	return template.New("strest").Funcs(tmpFuncs())
}

func tmpFuncs() template.FuncMap {
	return template.FuncMap{
		"ENV":                      getEnv,
		"MaleFirstName":            fake.MaleFirstName,
		"FemaleFirstName":          fake.FemaleFirstName,
		"FirstName":                fake.FirstName,
		"MaleLastName":             fake.MaleLastName,
		"FemaleLastName":           fake.FemaleLastName,
		"LastName":                 fake.LastName,
		"MalePatronymic":           fake.MalePatronymic,
		"FemalePatronymic":         fake.FemalePatronymic,
		"Patronymic":               fake.Patronymic,
		"MaleFullNameWithPrefix":   fake.MaleFullNameWithPrefix,
		"FemaleFullNameWithPrefix": fake.FemaleFullNameWithPrefix,
		"FullNameWithPrefix":       fake.FullNameWithPrefix,
		"MaleFullNameWithSuffix":   fake.MaleFullNameWithSuffix,
		"FemaleFullNameWithSuffix": fake.FemaleFullNameWithSuffix,
		"FullNameWithSuffix":       fake.FullNameWithSuffix,
		"MaleFullName":             fake.MaleFullName,
		"FemaleFullName":           fake.FemaleFullName,
		"FullName":                 fake.FullName,
	}
}

func getEnv(in string) string {
	return os.Getenv(in)
}
