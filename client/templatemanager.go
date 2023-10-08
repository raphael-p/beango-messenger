package client

import (
	"html/template"
	"net/http"

	"github.com/raphael-p/beango/utils/logger"
	"github.com/raphael-p/beango/utils/response"
)

var templateMap map[string]*template.Template = make(map[string]*template.Template)

func getTemplate(name, value string) (*template.Template, error) {
	templateFromMap := templateMap[name]
	if templateFromMap != nil {
		return templateFromMap, nil
	}

	newTemplate, err := template.New(name).Parse(value)
	if err != nil {
		return nil, err
	}
	templateMap[name] = newTemplate
	return newTemplate, nil
}

func ServeTemplate(w *response.Writer, name, value string, data map[string]any) {
	newTemplate, err := getTemplate(name, value)
	if err != nil {
		logger.Error(err.Error())
		w.WriteString(http.StatusInternalServerError, err.Error())
		return
	}
	if err := newTemplate.Execute(w, data); err != nil {
		logger.Error(err.Error())
		w.WriteString(http.StatusInternalServerError, err.Error())
		return
	}
}
