package main

import (
	"bytes"
	"fmt"
	"internal/adapters/logger"
	"internal/domain"
	"text/template"

	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
	//check for malformed requests - only exact root path accepted
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// set correct data type
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	qr, err := stor.GetQuizRows(false)
	if err != nil {
		logger.Error(fmt.Sprintf("Error getting quiz rows: %v\n", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := domain.PageData{
		Title:       "Прохождение тестирования",
		TableHeader: "Доступные тесты:",
		Columns: []string{
			"Название",
			"Описание",
			"Тест",
		},
		Rows: *qr,
	}

	//TODO: move it and cache it!
	// Read template
	tmpl, err := template.ParseFiles("templates\\index.html")
	logger.Info(fmt.Sprintf("ParseFiles => name is: %s %s\n", tmpl.Name(), tmpl.DefinedTemplates()))

	if err != nil {
		logger.Error(fmt.Sprintf("Error loading template: %v\n", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//genereate from string
	// t := template.Must(template.New("webpage").Parse(tmpl))

	buf := &bytes.Buffer{}
	//err = tmpl.Execute(buf, data)
	err = tmpl.ExecuteTemplate(buf, "page", data)
	if err != nil {
		logger.Error(fmt.Sprintf("Error executing template: %v\n", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// No error, send the content, HTTP 200 response status implied
	buf.WriteTo(w)
}
