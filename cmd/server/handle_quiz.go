package main

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"internal/adapters/logger"
	"internal/domain"

	"net/http"

	"github.com/go-chi/chi/v5"
)

func handleQuizPage(w http.ResponseWriter, r *http.Request) {

	// set correct data type
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	//w.WriteHeader(http.StatusOK)

	url := chi.URLParam(r, "id")
	uuid, err := domain.DecodeGUID(url)
	if err != nil {
		logger.Error(fmt.Sprintf("Wrong Quiz ID [%s]: %s\n", url, err.Error()))
		http.Error(w, fmt.Sprintf("Error. Resource not found: [%s]", url), http.StatusNotFound)
		return
	}

	if strings.TrimSpace(uuid) == "" {
		logger.Error(fmt.Sprintf("Malformed Quiz ID [%s]\n", url))
		http.Error(w, fmt.Sprintf("Error. Resource not found: [%s]", url), http.StatusNotFound)
		return
	}

	qd, err := stor.GetQuizData(uuid)
	if err != nil {
		logger.Error(fmt.Sprintf("Quiz extract error [%s]: %s\n", url, err.Error()))
		http.Error(w, fmt.Sprintf("Error. Resource not found: [%s]", url), http.StatusNotFound)
		return
	}

	if !qd.IsActive {
		logger.Error(fmt.Sprintf("Quiz is not active [%s]\n", url))
		http.Error(w, fmt.Sprintf("Error. Resource not found: [%s]", url), http.StatusNotFound)
		return
	}

	data := domain.QuizPageData{
		ID:          qd.ID,
		UUID:        qd.UUID,
		Title:       qd.Name,
		Description: qd.Description,
		Questions:   qd.Questions,
	}

	//TODO: move it and cache it!
	// Read template
	tmpl, err := template.ParseFiles("templates\\quiz.html")
	if err != nil {
		logger.Error(fmt.Sprintf("Error loading template: %v\n", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Info(fmt.Sprintf("ParseFiles => name is: %s %s", tmpl.Name(), tmpl.DefinedTemplates()))

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
