package main

import (
	"bytes"
	"fmt"
	"internal/adapters/logger"
	"internal/domain"
	"text/template"

	"net/http"
)

func adminPage(w http.ResponseWriter, r *http.Request) {

	if !isAuthenticated(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// set correct data type
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	//w.WriteHeader(http.StatusOK)

	qr, err := stor.GetQuizRows(true)
	if err != nil {
		logger.Error(fmt.Sprintf("Error getting quiz rows: %v\n", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := domain.PageData{
		Title:       "Панель администратора",
		TableHeader: "Доступные тесты:",
		Columns: []string{
			"UUID",
			"Идентификатор",
			"Название",
			"Описание",
			"Активен",
			"Тест",
			"Управление",
		},
		// Rows: []domain.QuizRowData{
		// 	{Name: "Тест 1", Value: "Отключить", Link: "#1"},
		// 	{Name: "Тест 2", Value: "Отключить", Link: "#2"},
		// 	{Name: "Тест 3", Value: "Отключить", Link: "#3"},
		// },
		Rows: *qr,
	}

	//TODO: move it and cache it!
	// Read template
	tmpl, err := template.ParseFiles("templates\\admin.html")
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
