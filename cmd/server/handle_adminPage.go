package main

import (
	"bytes"
	"fmt"
	"internal/adapters/logger"
	"text/template"

	"net/http"
)

// // Page data storage
// type PageData struct {
// 	Title       string
// 	TableHeader string
// 	Columns     []string
// 	Rows        []RowData
// }

// // Row description
// type RowData struct {
// 	Name  string
// 	Value string
// 	Link  string
// }

//const indexTableRowTpl = "<tr><td>%s</td><td style=\"text-align: right;\">%v</td></tr>"

func adminPage(w http.ResponseWriter, r *http.Request) {
	//check for malformed requests - only exact root path accepted
	//Important: covered by tests, removal will bring tests to fail
	// if r.URL.Path != "/" {
	// 	http.NotFound(w, r)
	// 	return
	// }

	// set correct data type
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	data := PageData{
		Title:       "Интерфейс администратора",
		TableHeader: "Доступные тесты:",
		Columns: []string{
			"Описание",
			"Тест",
		},
		Rows: []RowData{
			{Name: "Тест 1", Value: "Отключить", Link: "#1"},
			{Name: "Тест 2", Value: "Отключить", Link: "#2"},
			{Name: "Тест 3", Value: "Отключить", Link: "#3"},
		},
	}

	//TODO: move it and cache it!
	// Read template
	tmpl, err := template.ParseFiles("templates\\admin.html")
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
