package main

import (
	"fmt"
	"internal/adapters/logger"
	"internal/app"
	"net/http"
)

func submit(w http.ResponseWriter, r *http.Request) {
	//check for malformed requests - only exact root path accepted
	//Important: covered by tests, removal will bring tests to fail
	// if r.URL.Path != "/" {
	// 	http.NotFound(w, r)
	// 	return
	// }

	r.ParseForm()

	//DEBUG
	for key, values := range r.Form {
		logger.Info(fmt.Sprintf("TEST: %v\n", key))
		for _, value := range values {
			logger.Info(fmt.Sprintf(" --> Values: [%v][%v]\n", key, value))
		}
	}

	uuid := r.Form.Get("test_id")
	if uuid == "" {
		logger.Error(fmt.Sprintf("Submit: Wrong test ID [%s]\n", uuid))
		http.Error(w, "Error. Unexpected quiz ID", http.StatusBadRequest)
		return
	}

	qd, err := stor.GetQuizData(uuid)
	if err != nil {
		logger.Error(fmt.Sprintf("Quiz extract error [%s]: %s\n", uuid, err.Error()))
		http.Error(w, fmt.Sprintf("Error. Quiz not found [%s]", uuid), http.StatusInternalServerError)
		return
	}

	total_corrrect := app.GetQuizTotalScore(qd)
	total_answered := app.GetQuizUserScore(qd, r.Form)

	if total_corrrect == 0 {
		logger.Error(fmt.Sprintf("Zero correct answers in quiz [%s]\n", uuid))
		http.Error(w, fmt.Sprintf("Error. Quiz result calculation error [%s]", uuid), http.StatusInternalServerError)
		return
	}

	total_percent := app.GetQuizUserRatio(total_answered, total_corrrect)

	logger.Info(fmt.Sprintf("RESULTS: %2d/%2d (ratio %d%%)\n", total_answered, total_corrrect, total_percent))

	http.Redirect(w, r, "/results/1", http.StatusSeeOther)

}
