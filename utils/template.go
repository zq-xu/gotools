package utils

import (
	"bytes"
	"text/template"
)

func TransTemplate(tpl string, data interface{}) []byte {
	buffer := &bytes.Buffer{}
	_ = template.Must(template.New("").Parse(tpl)).Execute(buffer, data)
	return buffer.Bytes()
}
