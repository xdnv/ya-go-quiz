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

// Quiz question option
type quizOption struct {
	Value string
	Text  string
}

// Quiz question entry
type quizQuestion struct {
	Text    string
	Type    string // "single_choice", "multiple_choice" or "entry_int"
	Options []quizOption
}

// Page data structure
type quizPageData struct {
	Code        string
	Title       string
	Description string
	Questions   []quizQuestion
}

//const indexTableRowTpl = "<tr><td>%s</td><td style=\"text-align: right;\">%v</td></tr>"

func quiz(w http.ResponseWriter, r *http.Request) {
	//check for malformed requests - only exact root path accepted
	//Important: covered by tests, removal will bring tests to fail
	// if r.URL.Path != "/" {
	// 	http.NotFound(w, r)
	// 	return
	// }

	// set correct data type
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	logger.Info("alive 1")

	questions := []quizQuestion{
		{
			Text: "Выберите один правильный ответ (вариант №2)",
			Type: "single_choice",
			Options: []quizOption{
				{Value: "TST000001_001_01", Text: "Неправильный ответ №1"},
				{Value: "TST000001_001_02", Text: "Правильный ответ №2"},
				{Value: "TST000001_001_03", Text: "Неправильный ответ №3"},
			},
		},
		{
			Text: "Выберите все правильные ответы (варианты 2 и 4)",
			Type: "multiple_choice",
			Options: []quizOption{
				{Value: "TST000001_002_01", Text: "Неправильный ответ №1"},
				{Value: "TST000001_002_02", Text: "Правильный ответ №2"},
				{Value: "TST000001_002_03", Text: "Неправильный ответ №3"},
				{Value: "TST000001_002_04", Text: "Правильный ответ №4"},
			},
		},
		{
			Text: "Введите правильное значение (9)",
			Type: "entry_int",
			Options: []quizOption{
				{Value: "", Text: "(Введите число)"},
			},
		},
	}

	data := quizPageData{
		Code:        "000001",
		Title:       "Первый тест",
		Description: "Описание теста",
		Questions:   questions,
	}

	logger.Info("alive 2")

	//TODO: move it and cache it!
	// Read template
	tmpl, err := template.ParseFiles("templates\\quiz.html")
	if err != nil {
		logger.Error(fmt.Sprintf("Error loading template: %v\n", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Info(fmt.Sprintf("ParseFiles => name is: %s %s\n", tmpl.Name(), tmpl.DefinedTemplates()))

	logger.Info("alive 3")

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

	logger.Info("alive 4")

	// No error, send the content, HTTP 200 response status implied
	buf.WriteTo(w)
}
