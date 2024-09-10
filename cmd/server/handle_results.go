package main

import (
	"bytes"
	"fmt"
	"internal/adapters/logger"
	"internal/domain"
	"strings"
	"text/template"

	"net/http"

	"github.com/go-chi/chi/v5"
)

func handleResults(w http.ResponseWriter, r *http.Request) {

	// set correct data type
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	//w.WriteHeader(http.StatusOK)

	url := chi.URLParam(r, "id")
	uuid, err := domain.DecodeGUID(url)
	if err != nil {
		logger.Error(fmt.Sprintf("Wrong Result ID [%s]: %s\n", url, err.Error()))
		http.Error(w, fmt.Sprintf("Error. Resource not found: [%s]", url), http.StatusNotFound)
		return
	}

	if strings.TrimSpace(uuid) == "" {
		logger.Error(fmt.Sprintf("Malformed Result ID [%s]\n", url))
		http.Error(w, fmt.Sprintf("Error. Resource not found: [%s]", url), http.StatusNotFound)
		return
	}

	qr, err := stor.GetQuizResult(uuid)
	if err != nil {
		logger.Error(fmt.Sprintf("Result extract error [%s]: %s\n", url, err.Error()))
		http.Error(w, fmt.Sprintf("Error. Resource not found: [%s]", url), http.StatusNotFound)
		return
	}

	data := domain.QuizResultPageData{
		Title:  "Результаты тестирования",
		Link:   "https://example.com/results/" + domain.EncodeGUID(qr.TestID),
		Result: *qr,
	}

	//({{ if .IsPassed }}Верно{{ else }}Неверно{{ end }})

	//TODO: move it and cache it!
	// Read template
	tmpl, err := template.ParseFiles("templates\\results.html")
	logger.Info(fmt.Sprintf("ParseFiles => name is: %s %s", tmpl.Name(), tmpl.DefinedTemplates()))

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
