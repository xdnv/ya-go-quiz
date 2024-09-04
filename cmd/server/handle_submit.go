package main

import (
	"fmt"
	"internal/adapters/logger"
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
	for key, values := range r.Form {
		logger.Info(fmt.Sprintf("TEST: %v\n", key))
		for _, value := range values {
			logger.Info(fmt.Sprintf(" --> Values: [%v][%v]\n", key, value))
		}
	}
	http.Redirect(w, r, "/results/1", http.StatusSeeOther)

	// // Read body
	// body, err := io.ReadAll(r.Body)
	// if err != nil {
	// 	logger.Error(fmt.Sprintf("Error reading body: %v\n", err))
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	// defer r.Body.Close()

	// logger.Info(fmt.Sprintf("Got data: %s\n", body))

	// set correct data type
	//w.Header().Set("Content-Type", "text/html; charset=utf-8")
	//w.WriteHeader(http.StatusOK)
	//return
}
