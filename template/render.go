package template

import (
	"github.com/alexkappa/mustache"
)

func render(template string, context interface{}) (string, error) {
	t := mustache.New()

	err := t.ParseString(template)
	if err != nil {
		return "", err
	}

	result, err := t.RenderString(context)
	if err != nil {
		return "", err
	}

	return result, nil
}

func RenderIndex(context interface{}) (string, error) {
	return render(indexTemplate, context)
}
