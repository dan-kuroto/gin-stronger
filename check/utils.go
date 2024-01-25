package check

import (
	"bytes"
	"text/template"
)

func execErrorTemplate[T any](errTpl string, name string, value T) string {
	tpl, err := template.New("").Parse(errTpl)
	if err != nil {
		panic(err.Error())
	}

	var buffer bytes.Buffer
	err = tpl.Execute(&buffer, map[string]any{
		"name":  name,
		"value": value,
	})
	if err != nil {
		panic(err.Error())
	}

	return buffer.String()
}
